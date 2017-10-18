package algorithm

import (
	"fmt"
	"k8s.io/api/core/v1"
	"crypto/rsa"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
	"golang.org/x/oauth2/jws"
	"cit-extended-scheduler/util"
	"io/ioutil"
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

//ValidateCipherByPublicKey is used for validate the annotation(cipher) by public key
func ValidateCipherByPublicKey(cipherText string, key *rsa.PublicKey) error {
	validationStatus := jws.Verify(cipherText, key)
	return validationStatus
}

//JWTParseWithClaims is used for parsing and adding the annotation values in claims map
func JWTParseWithClaims(cipherText string, verifyKey *rsa.PublicKey, claim jwt.MapClaims) {
	_, err := jwt.ParseWithClaims(cipherText, claim, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
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
	fmt.Println("verifyKey:")
	fmt.Println(verifyKey)
	fmt.Println("trust annotation:")
	fmt.Println(cipherText)
	validationStatus := ValidateCipherByPublicKey(cipherText, verifyKey)
	fmt.Println("validation status", validationStatus)
	fmt.Println("claims before", claims)
	//cipherText is the annotation applied to the node, claims is the parsed AH report assigned as the annotation
	JWTParseWithClaims(cipherText, verifyKey, claims)
	fmt.Println("claims after", claims)

	trustedFlag, trustedVerifyFlag, locationFlag, locationVerifyFlag := ValidatePodWithAnnotation(node, claims)
	trustTimeFlag, assetTimeFlag := ValidateNodeByTime(claims)

	if validationStatus == nil && trustTimeFlag != 1 && assetTimeFlag != 1 {
		glog.V(4).Infof("Signature is valid but node valid time is expired")
		return false
	} else if (validationStatus == nil) && (trustTimeFlag == 1 && assetTimeFlag == 1) && (trustedFlag == 1 && trustedVerifyFlag == 1) && (locationVerifyFlag == 1 && locationFlag == 1) {
		glog.V(4).Infof("Signature is valid for both trusted and lacation tags ")
		return true
	} else if (validationStatus == nil) && (trustTimeFlag == 1) && (trustedFlag == 1 && trustedVerifyFlag == 1) {
		glog.V(4).Infof("Signature is valid for trusted tag")
		return true
	} else if (validationStatus == nil) && (assetTimeFlag == 1) && (locationVerifyFlag == 1 && locationFlag == 1) {
		glog.V(4).Infof("Signature is valid for location tag")
		return true
	} else {
		glog.Errorf("Signature validation failed")
		return false
	}
}
