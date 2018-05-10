/*
Copyright © 2018 Intel Corporation
SPDX-License-Identifier: BSD-3-Clause
*/

package api

import (
	"k8s_scheduler_cit_extension-k8s_extended_scheduler/algorithm"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"io/ioutil"
	schedulerapi "k8s.io/kubernetes/plugin/pkg/scheduler/api"
)

var Confpath string

type Config struct {
	Trusted string `"json":"trusted"`
}

func getPrefixFromConf(path string) (string, error) {
        out, err := ioutil.ReadFile(path)
        if err != nil {
                glog.Errorf("Error: %s %v", path, err)
                return "",err
        }
        s := Config{}
        err = json.Unmarshal(out, &s)
        if err != nil {
                glog.Errorf("Error:  %v", err)
                return "",err
        }
        return s.Trusted, nil
}

//FilterHandler is the filter host.
func FilterHandler(c *gin.Context) {
	var args schedulerapi.ExtenderArgs
	glog.V(4).Infof("Post received at extended scheduler: %v", args)
	//Create a binding for args passed to the POST api
	if c.BindJSON(&args) == nil {
		prefixString,er := getPrefixFromConf(Confpath)
		if er != nil {
			glog.Fatalf("Error:%v",er)
		}
		result, err := algorithm.FilteredHost(&args, prefixString)
		if err == nil {
			c.JSON(200, result)
		} else {
			c.JSON(500, err)
		}
	}
}
