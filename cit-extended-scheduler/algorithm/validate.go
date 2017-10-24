package algorithm

import (
	"fmt"
	"regexp"
	"k8s.io/api/core/v1"
	"strconv"
	"strings"
	"time"
	jwt "github.com/dgrijalva/jwt-go"
)

const (
	ahreport string = "AssetTagSignedReport"
)


//ValidatePodWithAnnotation is to validate signed trusted and location report with pod keys and values
func ValidatePodWithAnnotation(podData []v1.NodeSelectorRequirement, claims jwt.MapClaims) (int, int, int, int) {
	trustFlag, trustedVerifyFlag, locationFlag, locationVerifyFlag, annotateLocationFlag := 0, 0, 0, 0, 0

	for _, val := range podData {
		if nodeVal, ok := claims[val.Key]; ok {
			trustFlag = 0
			trustedVerifyFlag = 1
			for _, podVal := range val.Values {
				if nodeVal == true || nodeVal == false {
					nodeValTemp := nodeVal.(bool)
					nodeVal := strconv.FormatBool(nodeValTemp)
					if podVal == nodeVal {
						trustFlag = 1
						break
					}
				}
			}
		} else {
			if strings.Contains(val.Key, ".") {
				podValArray := strings.Split(val.Key, ".")
				if claims[ahreport].(map[string]interface{})[podValArray[0]] != nil {
					locationVerifyFlag = 1
					assetTagMap := claims[ahreport].(map[string]interface{})[podValArray[0]]
					assetTagList, ok := assetTagMap.([]interface{})
					annotateLocationFlag = 0
					if ok {
						for _, assetTagValue := range assetTagList {
							if podValArray[1] == assetTagValue {
								annotateLocationFlag = 1
								break
							}
						}
					}
					if annotateLocationFlag == 0 {
						locationFlag = 0
					} else {
						locationFlag = 1
					}
				}
			}
		}
	}
	fmt.Println("trustFlag, trustedVerifyFlag, locationFlag, locationVerifyFlag")
	fmt.Println(trustFlag, trustedVerifyFlag, locationFlag, locationVerifyFlag)
	return trustFlag, trustedVerifyFlag, locationFlag, locationVerifyFlag
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
		fmt.Println(newstr)
		trustedValidToTime := strings.Replace(timeVal, newstr, "", -1)
		fmt.Println("Trust validity time", timeVal )
		fmt.Println("Trust validity time after replace ", trustedValidToTime )

		t := time.Now()
		timeDiff := strings.Compare(trustedValidToTime, t.Format(time.RFC3339))
		fmt.Println("Time Now:", t.Format(time.RFC3339))
		fmt.Println("Time diff:", timeDiff)
		if timeDiff >= 0 {
			trustedTimeFlag = 1
		}
	}

	return trustedTimeFlag
}
