/*
Copyright © 2018 Intel Corporation
SPDX-License-Identifier: BSD-3-Clause
*/

package algorithm

import (
	"fmt"
	"k8s_scheduler_cit_extension-k8s_extended_scheduler/util"
	"crypto/rsa"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
	"golang.org/x/oauth2/jws"
	"k8s.io/api/core/v1"
)

//fatal functions just logs and exits
func fatal(err error) {
	if err != nil {
		//panic(err.Error())
		fmt.Println(err.Error())
	}
}

//ParseRSAPublicKeyFromPEM is used for parsing and verify public key
func ParseRSAPublicKeyFromPEM(pubKey []byte) (*rsa.PublicKey, error) {
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKey)
	if err != nil {
		fmt.Println("error in ParseRSAPublicKeyFromPEM")
	}
	return verifyKey, err
}

//ValidateAnnotationByPublicKey is used for validate the annotation(cipher) by public key
func ValidateAnnotationByPublicKey(cipherText string, key *rsa.PublicKey) error {
	validationStatus := jws.Verify(cipherText, key)
	return validationStatus
}

//JWTParseWithClaims is used for parsing and adding the annotation values in claims map
func JWTParseWithClaims(cipherText string, verifyKey *rsa.PublicKey, claim jwt.MapClaims) {
	token, err := jwt.ParseWithClaims(cipherText, claim, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	glog.V(4).Infof("Parsed token is :", token)
	if err != nil {
		fmt.Println("error in JWTParseWithClaims")
		fatal(err)
	}
}

//CheckAnnotationAttrib is used for validate node with restpect to time,trusted and location tags
func CheckAnnotationAttrib(cipherText string, node []v1.NodeSelectorRequirement, trustPrefix string) bool {
	var claims = jwt.MapClaims{}
	pubKey := util.GetAHPublicKey()
	verifyKey, err := ParseRSAPublicKeyFromPEM(pubKey)
	if err != nil {
		fatal(err)
		glog.V(4).Infof("Invalid AH public key")
		return false
	}
	validationStatus := ValidateAnnotationByPublicKey(cipherText, verifyKey)
	//fmt.Println("validation status", validationStatus)
	//fmt.Println("claims before", claims)
	if validationStatus == nil {
		glog.V(4).Infof("Signature is valid, STR is from valid AH")
	} else {
		fatal(validationStatus)
		glog.Errorf("Signature validation failed")
		return false
	}

	//cipherText is the annotation applied to the node, claims is the parsed AH report assigned as the annotation
	JWTParseWithClaims(cipherText, verifyKey, claims)
	fmt.Println("claims after", claims)

	verify := ValidatePodWithAnnotation(node, claims, trustPrefix)
	if verify {
		glog.V(4).Infof("Node label validated against node annotations succesfull")
	} else {
		glog.Errorf("Node Label did not match node annotation ")
		return false
	}

	trustTimeFlag := ValidateNodeByTime(claims)

	if trustTimeFlag == 1 {
		glog.V(4).Infof("Attested node validity time check passed")
		return true
	} else {
		glog.Errorf("Attested node validity time has expired")
		return false
	}
}
