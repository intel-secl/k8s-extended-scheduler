/*
Copyright © 2018 Intel Corporation
SPDX-License-Identifier: BSD-3-Clause
*/

package util

import (
	"github.com/tkanos/gonfig"
	"strconv"
	"log"
)

var AH_KEY_FILE string

func GetCmdlineArgs() (string, string, string, string) {

	type extenedSchedConfig struct {
		Url  string //Extended scheduler url
		Port int    //Port for the Extended scheduler to listen on
		//Server Certificate to be used for TLS handshake
		ServerCert string
		//Server Key to be used for TLS handshake
		ServerKey string
		//Attestation Hub Key to be used for parsing signed trust report
		AttestationHubKey string
	}

	conf := extenedSchedConfig{}
	schedConf := "/opt/isecl-k8s-extensions/config/isecl-extended-scheduler-config.json"
	err := gonfig.GetConf(schedConf, &conf)
	if err != nil {
		log.Fatalf("Error: Please ensure extended schduler configuration is present in curent dir,%v",err)
	}

	//PORT for the extended scheduler to listen.
	port_no := conf.Port
	port := strconv.Itoa(port_no)

	AH_KEY_FILE = (conf.AttestationHubKey)
	return conf.Url, port, conf.ServerCert, conf.ServerKey
}
