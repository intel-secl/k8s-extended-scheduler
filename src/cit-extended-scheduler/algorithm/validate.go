package algorithm

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"k8s.io/api/core/v1"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	//ahreport string = "AssetTagSignedReport"
	ahreport string = "asset_tags"
)

//ValidatePodWithAnnotation is to validate signed trusted and location report with pod keys and values
func ValidatePodWithAnnotation(nodeData []v1.NodeSelectorRequirement, claims jwt.MapClaims) bool {
	//trustFlag, trustedVerifyFlag, locationFlag, locationVerifyFlag, annotateLocationFlag := 0, 0, 0, 0, 0
	//annotateLocationFlag := 0
	assetClaims := claims[ahreport].(map[string]interface{})
	fmt.Println("Asset tag report is ", assetClaims)

	for _, val := range nodeData {
		//if val is trusted, it can be directlly found in claims
		if sigVal, ok := claims[val.Key]; ok {
			for _, nodeVal := range val.Values {
				if sigVal == true || sigVal == false {
					fmt.Println("sigVal is boolean")
					sigValTemp := sigVal.(bool)
					sigVal := strconv.FormatBool(sigValTemp)
					if nodeVal == sigVal {
						fmt.Println("Trusted val found")
						continue
					} else {
						fmt.Println("Trust tag is tampered")
						return false
					}
				} else {
					fmt.Println("sigVal is not boolean")
					if nodeVal == sigVal {
						fmt.Println("Trusted val found")
						continue
					} else {
						fmt.Println("Trust tag is tampered")
						return false
					}
				}
			}
		} else {
			if strings.Contains(val.Key, ".") {
				fmt.Println("Attestation hub val ", val)
				nodeValArray := strings.Split(val.Key, ".")

				if geoKey, ok := assetClaims[nodeValArray[0]]; ok {
					fmt.Println("Found Key in AH report", geoKey)
					//assetTagMap := claims[ahreport].(map[string]interface{})[podValArray[0]]
					//assetTagList, ok := assetTagMap.([]interface{})
					assetTagList, ok := geoKey.([]interface{})
					if ok {
						for _, geoVal := range assetTagList {
							if nodeValArray[1] == geoVal {
								fmt.Println("Asset tag value found in AH report")
								continue
							} else {
								fmt.Println("Asset tag value not found in AH report")
								return false
							}
						}
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
		//trustedValidToTime = strings.Replace(timeVal, ".", ":", -1)
		reg, err := regexp.Compile("[0-9]+-[0-9]+-[0-9]+T[0-9]+:[0-9]+:[0-9]+")
		if err != nil {
			fmt.Println(err)
		}
		newstr := reg.ReplaceAllString(timeVal, "")
		//fmt.Println(newstr)
		trustedValidToTime := strings.Replace(timeVal, newstr, "", -1)
		//fmt.Println("Trust validity time", timeVal )
		//fmt.Println("Trust validity time after replace ", trustedValidToTime )

		t := time.Now().UTC()
		timeDiff := strings.Compare(trustedValidToTime, t.Format(time.RFC3339))
		//fmt.Println("Time Now:", t.Format(time.RFC3339))
		//fmt.Println("Time diff:", timeDiff)
		if timeDiff >= 0 {
			trustedTimeFlag = 1
		}
	}

	return trustedTimeFlag
}
