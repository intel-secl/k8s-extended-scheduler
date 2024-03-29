/*
Copyright © 2019 Intel Corporation
SPDX-License-Identifier: BSD-3-Clause
*/

package algorithm

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"strings"
	"k8s_scheduler_cit_extension-k8s_extended_scheduler/util"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
	"k8s.io/api/core/v1"
)

//ParseRSAPublicKeyFromPEM is used for parsing and verify public key
func ParseRSAPublicKeyFromPEM(pubKey []byte) (*rsa.PublicKey, error) {
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKey)
	if err != nil {
		glog.Errorf("error in ParseRSAPublicKeyFromPEM")
		return nil,err
	}
	return verifyKey, err
}

//ValidateAnnotationByPublicKey is used for validating the annotation(cipher) by public key
func ValidateAnnotationByPublicKey(cipherText string, key *rsa.PublicKey) error {
	parts := strings.Split(cipherText, ".")
	if len(parts) != 3 {
		return errors.New("jws: invalid token received, token must have 3 parts")
	}

	signedContent, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
                return err
        }

	signatureString, err := base64.StdEncoding.DecodeString(parts[2])
	if err != nil {
		return err
	}

	h := sha512.New384()
	h.Write(signedContent)
	return rsa.VerifyPKCS1v15(key, crypto.SHA384, h.Sum(nil), signatureString)
}

//JWTParseWithClaims is used for parsing and adding the annotation values in claims map
func JWTParseWithClaims(cipherText string, verifyKey *rsa.PublicKey, claim jwt.MapClaims) {
	token, err := jwt.ParseWithClaims(cipherText, claim, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	glog.Infof("Parsed token is :", token)
	if err != nil {
		glog.Errorf("error in JWTParseWithClaims")
	}
}

//CheckAnnotationAttrib is used to validate node with respect to time,trusted and location tags
func CheckAnnotationAttrib(cipherText string, node []v1.NodeSelectorRequirement, trustPrefix string) bool {
	var claims = jwt.MapClaims{}
	pubKey := util.GetAHPublicKey()
	verifyKey, err := ParseRSAPublicKeyFromPEM(pubKey)
	if err != nil {
		glog.Errorf("Invalid AH public key")
		return false
	}
	validationError := ValidateAnnotationByPublicKey(cipherText, verifyKey)
	if validationError == nil {
		glog.Infof("Signature is valid, STR is from valid AH")
	} else {
		glog.Errorf("Signature validation failed, Error: %v", validationError)
		return false
	}

	//cipherText is the annotation applied to the node, claims is the parsed AH report assigned as the annotation
	JWTParseWithClaims(cipherText, verifyKey, claims)

	glog.Infof("CheckAnnotationAttrib - Parsed claims for %v",  claims)

	verify := ValidatePodWithAnnotation(node, claims, trustPrefix)
	if verify {
		glog.Infoln("Node label validated against node annotations succesful")
	} else {
		glog.Infoln("Node Label did not match node annotation ")
		return false
	}

	trustTimeFlag := ValidateNodeByTime(claims)

	if trustTimeFlag == 1 {
		glog.Infoln("Attested node validity time check passed")
		return true
	} else {
		glog.Infoln("Attested node validity time has expired")
		return false
	}
}
