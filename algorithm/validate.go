/*
Copyright © 2018 Intel Corporation
SPDX-License-Identifier: BSD-3-Clause
*/

package algorithm

import (
	jwt "github.com/dgrijalva/jwt-go"
	"k8s.io/api/core/v1"
	"github.com/golang/glog"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	ahreport string = "asset_tags"
	trusted  string = "trusted"
)

//ValidatePodWithAnnotation is to validate signed trusted and location report with pod keys and values
func ValidatePodWithAnnotation(nodeData []v1.NodeSelectorRequirement, claims jwt.MapClaims, trustprefix string) bool {
	assetClaims := claims[ahreport].(map[string]interface{})

	for _, val := range nodeData {
		//if val is trusted, it can be directly found in claims
		if sigVal, ok := claims[trusted]; ok {
			tr := trustprefix + trusted
			if val.Key == tr {
				for _, nodeVal := range val.Values {
					if sigVal == true || sigVal == false {
						sigValTemp := sigVal.(bool)
						sigVal := strconv.FormatBool(sigValTemp)
						if nodeVal == sigVal {
							continue
						} else {
							return false
						}
					} else {
						if nodeVal == sigVal {
							continue
						} else {
							return false
						}
					}
				}
			}
		} else {
			if geoKey, ok := assetClaims[val.Key]; ok {
				assetTagList, ok := geoKey.([]interface{})
				if ok {
					flag := false
					//Taking only first value from asset tag list assuming only one value will be there
					geoVal := assetTagList[0]
					newVal := geoVal.(string)
					newVal = strings.Replace(newVal, " ", "", -1)
					for _, match := range val.Values {
						if match == newVal {
							flag = true
						}
					}
					if flag {
						continue
					} else {
						return false
					}
				}

			}
		}
	}
	return true
}

//ValidateNodeByTime is used for validate time for each node with current system time(Expiry validation)
func ValidateNodeByTime(claims jwt.MapClaims) int {
	trustedTimeFlag := 0
	if timeVal, ok := claims["valid_to"].(string); ok {
		reg, err := regexp.Compile("[0-9]+-[0-9]+-[0-9]+T[0-9]+:[0-9]+:[0-9]+")
		if err != nil {
			glog.Errorf("%v",err)
		}
		newstr := reg.ReplaceAllString(timeVal, "")
		trustedValidToTime := strings.Replace(timeVal, newstr, "", -1)

		t := time.Now().UTC()
		timeDiff := strings.Compare(trustedValidToTime, t.Format(time.RFC3339))
		if timeDiff >= 0 {
			trustedTimeFlag = 1
		}
	}

	return trustedTimeFlag
}
