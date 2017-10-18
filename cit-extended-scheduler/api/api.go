package api

import (
	"fmt"
	"cit-extended-scheduler/algorithm"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	schedulerapi "k8s.io/kubernetes/plugin/pkg/scheduler/api"
)

//FilterHandler is the filter host.
func FilterHandler(c *gin.Context) {
	var args schedulerapi.ExtenderArgs
	glog.V(4).Infof("Post received at extended scheduler: %v", args)
	fmt.Println("Post received at extended scheduler: %v", args)
	//Create a binding for args passed to the POST api
	c.BindJSON(&args)
	c.JSON(200, algorithm.FilteredHost(&args))
}
