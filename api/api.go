/*
Copyright © 2018 Intel Corporation
SPDX-License-Identifier: BSD-3-Clause
*/

package api

import (
	"k8s_scheduler_cit_extension-k8s_extended_scheduler/algorithm"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"io/ioutil"
	schedulerapi "k8s.io/kubernetes/plugin/pkg/scheduler/api"
)

type Config struct {
	Trusted string `"json":"trustedPrefix"`
}

const (
	CONFPATH string = "/opt/k8s_scheduler_cit_extension-k8s_extended_scheduler/bin/tag_prefix.conf"
)

func getPrefixFromConf(path string) string {
	out, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	s := Config{}
	err = json.Unmarshal(out, &s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Read from config ...................")
	fmt.Println(s.Trusted)
	return s.Trusted
}

//FilterHandler is the filter host.
func FilterHandler(c *gin.Context) {
	var args schedulerapi.ExtenderArgs
	glog.V(4).Infof("Post received at extended scheduler: %v", args)
	//fmt.Println("Post received at extended scheduler: %v", args)
	//Create a binding for args passed to the POST api
	if c.BindJSON(&args) == nil {
		prefixString := getPrefixFromConf(CONFPATH)
		result, err := algorithm.FilteredHost(&args, prefixString)
		if err == nil {
			c.JSON(200, result)
		} else {
			c.JSON(500, err)
		}
	}
}
