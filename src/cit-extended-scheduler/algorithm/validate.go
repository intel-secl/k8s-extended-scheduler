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
	ahreport string = "asset_tags"
)

//ValidatePodWithAnnotation is to validate signed trusted and location report with pod keys and values
func ValidatePodWithAnnotation(nodeData []v1.NodeSelectorRequirement, claims jwt.MapClaims) bool {
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
					assetTagList, ok := geoKey.([]interface{})
					if ok {
						fmt.Println("nodeVal value is ", nodeValArray[1])
						flag := false
						for _, geoVal := range assetTagList {
							newVal := geoVal.(string)
							newVal = strings.Replace(newVal, " ", "", -1)
							fmt.Println("geoVal value is ", newVal)
							if nodeValArray[1] == newVal {
								fmt.Println("Asset tag value found in AH report")
								flag = true
							}
						}
						if flag {
							fmt.Println("Asset tag value found in AH report flag is true")
							continue
						} else {
							fmt.Println("Asset tag not found flag is false")
							return false
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
