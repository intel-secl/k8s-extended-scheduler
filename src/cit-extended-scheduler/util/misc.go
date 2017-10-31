package util

import (
	"fmt"
	//"flag"
	"github.com/tkanos/gonfig"
	"strconv"
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
	//schedConf := flag.String("schedConf", "", "Configration file for Extended Scheduler")
	//flag.Parse()
	schedConf := "/opt/cit_k8s_extensions/bin/cit-extended-scheduler-config.json"
	//err := gonfig.GetConf("./extended_scheduler_config.json", &conf)
	/*
		if *schedConf == "" {
			fmt.Println("No Extended Scheduler configuration passed")
			panic("Needs Extended Scheduler Config")
		}
	*/
	//fmt.Println(schedConf)
	err := gonfig.GetConf(schedConf, &conf)
	if err != nil {
		fmt.Println("Error: Please ensure extended schduler configuration is present in curent dir")
		panic(err)
	}

	//PORT for the extended scheduler to listen.
	port_no := conf.Port
	port := strconv.Itoa(port_no)

	AH_KEY_FILE = (conf.AttestationHubKey)
	//fmt.Println( conf.Url, port, conf.ServerCert, conf.ServerKey)
	return conf.Url, port, conf.ServerCert, conf.ServerKey
}
