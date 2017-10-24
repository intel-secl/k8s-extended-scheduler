package algorithm

import (
	"fmt"
	//"strings"
	"k8s.io/api/core/v1"
	"crypto/rsa"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
	"golang.org/x/oauth2/jws"
	"cit-extended-scheduler/util"
	//"crypto/sha256"
	//"encoding/base64"
	//"io/ioutil"
)


//fatal functions just logs and exits
func fatal(err error) {
	if err != nil {
		//panic(err.Error())
		fmt.Println(err.Error())
	}
}

//ParseRSAPublicKeyFromPEM is used for parsing and verify public key
func ParseRSAPublicKeyFromPEM(pubKey []byte) *rsa.PublicKey {
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKey)
	if err != nil {
		fmt.Println("error in ParseRSAPublicKeyFromPEM")
	}
	fatal(err)
	return verifyKey
}

//ValidateAnnotationByPublicKey is used for validate the annotation(cipher) by public key
func ValidateAnnotationByPublicKey(cipherText string, key *rsa.PublicKey) error {
	validationStatus := jws.Verify(cipherText, key)
	return validationStatus
}

//JWTParseWithClaims is used for parsing and adding the annotation values in claims map
func JWTParseWithClaims(cipherText string, verifyKey *rsa.PublicKey, claim jwt.MapClaims) {
	token, err := jwt.ParseWithClaims(cipherText, claim, func(token *jwt.Token) ( interface{}, error) {
		return verifyKey, nil
	})
	fmt.Println("ganesh token is ", token)
	if err != nil {
		fmt.Println("error in JWTParseWithClaims")
	}
	fatal(err)
}

//CheckAnnotationAttrib is used for validate node with restpect to time,trusted and location tags
func CheckAnnotationAttrib(cipherText string, node []v1.NodeSelectorRequirement) bool {
	var claims = jwt.MapClaims{}
	pubKey := util.GetAHPublicKey()
	verifyKey := ParseRSAPublicKeyFromPEM(pubKey)
	validationStatus := ValidateAnnotationByPublicKey(cipherText, verifyKey)
	fmt.Println("validation status", validationStatus)
	fmt.Println("claims before", claims)
	//cipherText is the annotation applied to the node, claims is the parsed AH report assigned as the annotation
	JWTParseWithClaims(cipherText, verifyKey, claims)
	fmt.Println("claims after", claims)
	
	//trustedFlag, trustedVerifyFlag, locationFlag, locationVerifyFlag := ValidatePodWithAnnotation(node, claims)
	trustTimeFlag := ValidateNodeByTime(claims)

	if validationStatus == nil && trustTimeFlag != 1 {
		glog.V(4).Infof("Signature is valid but node validity time has expired")
		fmt.Println("Signature is valid but node valid time is expired")
		return false
	} else {
		glog.Errorf("Signature validation failed")
		fmt.Println("Signature validation failed")
		return false
	}
}
