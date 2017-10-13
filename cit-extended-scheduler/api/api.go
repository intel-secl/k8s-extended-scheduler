package api

import (
	"cit-extended-scheduler/algorithm"
	"github.com/gin-gonic/gin"
	schedulerapi "k8s.io/kubernetes/plugin/pkg/scheduler/api"
)

//FilterHandler is the filter host.
func FilterHandler(c *gin.Context) {
	var args schedulerapi.ExtenderArgs
	c.BindJSON(&args)
	c.JSON(200, algorithm.FilteredHost(&args))
}
