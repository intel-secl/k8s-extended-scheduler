package algorithm

import (
	"fmt"
	"k8s.io/api/core/v1"
	schedulerapi "k8s.io/kubernetes/plugin/pkg/scheduler/api"
        "io/ioutil"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2/jws"
	"log"
	"crypto/rsa"
	//"reflect"
	"strings"
	"time"
	"strconv"
)


//FilteredHost is used for getting the nodes and pod details and verify and return if pod key matches with annotations
func FilteredHost(args *schedulerapi.ExtenderArgs)  (*schedulerapi.ExtenderFilterResult)  {
	result := []v1.Node{}
	failedNodesMap := schedulerapi.FailedNodesMap{}
	nodes:= args.Nodes
	pod:= args.Pod

	if pod.Spec.Affinity != nil && pod.Spec.Affinity.NodeAffinity != nil && pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution != nil && pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms != nil{
		nodeSelectorData := pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms
	
		for _, node := range nodes.Items {
			if cipherVal, ok := node.Annotations["TrustTagSignedReport"]; ok {
				for _, nodeSelector := range nodeSelectorData {
					if(CheckAnnotationAttrib(cipherVal,nodeSelector.MatchExpressions)){
						result = append(result, node)
					}else{
						failedNodesMap[node.Name] = fmt.Sprintf("Annotation validation failed in extended-scheduler")					
					}
				}	
			}
		}
	}
	
	return &schedulerapi.ExtenderFilterResult{
		Nodes:       &v1.NodeList{Items: result},
		NodeNames:   nil,
		FailedNodes: failedNodesMap,
	}
}


var claims = jwt.MapClaims{}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//ReadPublicKey is used for reading and return the public key from particular file location 
func ReadPublicKey() []byte{
	pubKey,err := ioutil.ReadFile("attestaton_hub_keys/attestaton_hub_keys/hub_public_key.pem")
	fatal(err)
	return pubKey
}

//ParseRSAPublicKeyFromPEM is used for parsing and verify public key
func ParseRSAPublicKeyFromPEM(pubKey []byte) *rsa.PublicKey{
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKey)
	fatal(err)
	return verifyKey
}

//ValidateCipherByPublicKey is used for validate the annotation(cipher) by public key
func ValidateCipherByPublicKey(cipherText string,key *rsa.PublicKey) error{
	validationStatus := jws.Verify(cipherText,key)
	return validationStatus
}

//JWTParseWithClaims is used for parsing and adding the annotation values in claims map
func JWTParseWithClaims(cipherText string,verifyKey *rsa.PublicKey){
	_, err := jwt.ParseWithClaims(cipherText, claims, func(token *jwt.Token) (interface{}, error) {
	    return verifyKey, nil
	})
	fatal(err)
}

//CheckAnnotationAttrib is used for validate node with restpect to time,trusted and location tags
func CheckAnnotationAttrib(cipherText string,node []v1.NodeSelectorRequirement) bool{
	pubKey := ReadPublicKey()
	verifyKey := ParseRSAPublicKeyFromPEM(pubKey)
	validationStatus := ValidateCipherByPublicKey(cipherText,verifyKey)
	JWTParseWithClaims(cipherText,verifyKey)
	
	trustedFlag,trustedVerifyFlag,locationFlag,locationVerifyFlag := ValidatePodWithAnnotation(node)
	trustTimeFlag,assetTimeFlag := ValidateNodeByTime()

	if validationStatus == nil && trustTimeFlag != 1 && assetTimeFlag != 1{
		fmt.Println("Signature is valid but node valid time is expired")
		return false
	}else if (validationStatus == nil) && (trustTimeFlag == 1 && assetTimeFlag == 1) && (trustedFlag == 1 && trustedVerifyFlag == 1) && (locationVerifyFlag == 1 && locationFlag == 1){
		fmt.Println("Signature is valid by both trusted and lacation tags")
		return true
	}else if (validationStatus == nil) && (trustTimeFlag == 1) && (trustedFlag == 1 && trustedVerifyFlag == 1){
		fmt.Println("Signature is valid by trusted tag")
		return true
	}else if (validationStatus == nil) && (assetTimeFlag == 1) && (locationVerifyFlag == 1 && locationFlag == 1){
		fmt.Println("Signature is valid by lacation tag")
		return true
	}else {
		fmt.Println("Signature is not valid")
		return false
	}
}

//ValidatePodWithAnnotation is to validate signed trusted and location report with pod keys and values
func ValidatePodWithAnnotation(pod []v1.NodeSelectorRequirement) (int,int,int,int){
	trustFlag,trustedFlag,trustedVerifyFlag,locationFlag,locationVerifyFlag,annotateLocationFlag := 0,0,0,0,0,0
	for _, val := range pod {
		if nodeVal, ok := claims[val.Key]; ok {
			trustFlag = 0
			trustedVerifyFlag = 1
			for _, podVal := range val.Values {
				if nodeVal == true || nodeVal == false {
					nodeValTemp := nodeVal.(bool)
					nodeVal := strconv.FormatBool(nodeValTemp)
					if(podVal == nodeVal){
						trustFlag = 1
						break
					}
				}
			}
			if trustFlag == 0 {
				trustedFlag = 0
			} else{
				trustedFlag = 1
			}			
		}else{
			if(strings.Contains(val.Key, ".")){
				podValArray := strings.Split(val.Key, ".")
				if(claims["AssetTagSignedReport"].(map[string]interface{})[podValArray[0]] != nil){
					locationVerifyFlag = 1
					assetTagMap := claims["AssetTagSignedReport"].(map[string]interface{})[podValArray[0]]
					assetTagList, ok := assetTagMap.([]interface{})
					annotateLocationFlag = 0
					if ok {
						for _, assetTagValue := range assetTagList {
							if(podValArray[1] == assetTagValue){
								annotateLocationFlag = 1
								break
							}
						}	
	    				}
					if annotateLocationFlag == 0 {
						locationFlag = 0
					} else{
						locationFlag = 1
					}
				}		
			}
		}
	}
	return trustedFlag,trustedVerifyFlag,locationFlag,locationVerifyFlag
}

//ValidateNodeByTime is used for validate time for each node with current system time(Expiry validation)
func ValidateNodeByTime() (int,int){
	trustedValidToTime,assetValidToTime := "",""
	trustedTimeFlag,assetTimeFlag := 0,0
	if timeVal, ok := claims["TrustTagExpiry"].(string); ok {
		trustedValidToTime = strings.Replace(timeVal, ".",":",-1)
	}
	
	if timeVal, ok := claims["AssetTagExpiry"].(string); ok {
		assetValidToTime = strings.Replace(timeVal, ".",":",-1)
	}
	
	t := time.Now()
	if trustedTimeFlag == 0{
		timeDiff := Compare(trustedValidToTime,t.Format(time.RFC3339))
		if timeDiff >=0 {
			trustedTimeFlag = 1
		}
	}
	if assetTimeFlag == 0{
		timeDiff := Compare(assetValidToTime,t.Format(time.RFC3339))
		if timeDiff >=0 {
			assetTimeFlag = 1
		}
	}
	return trustedTimeFlag,assetTimeFlag
}


//Compare for compare to string
func Compare(a, b string) int {
  	if a == b {
  		return 0
  	}
  	if a < b {
  		return -1
  	}
  	return +1
}
