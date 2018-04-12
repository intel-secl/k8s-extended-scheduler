/*
Copyright © 2018 Intel Corporation
SPDX-License-Identifier: BSD-3-Clause
*/

package main

import (
	"k8s_scheduler_cit_extension-k8s_extended_scheduler/api"
	"k8s_scheduler_cit_extension-k8s_extended_scheduler/util"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"net/http"
	"os"
	"os/signal"
	"time"
	"flag"
)


func extendedScheduler(c *gin.Context) {
	c.JSON(200, gin.H{"result": "Cit Extended Scheduler"})
	return
}

func SetupRouter() (*gin.Engine, *http.Server) {
	//get a webserver instance, that contains a muxer, middleware and configuration settings
	router := gin.Default()

	// fetch all the cmd line args
	url, port, server_crt, server_key := util.GetCmdlineArgs()
	fmt.Println(server_crt, server_key)

	//initialize http server config
	server := &http.Server{
		Addr:    url + ":" + port,
		Handler: router,
	}

	//run the server instance
	go func() {
		// service connections
		if err := server.ListenAndServeTLS(server_crt, server_key); err != nil {
			glog.V(4).Infof("listen: %s\n", err)
			//fmt.Printf("listen %s ...", err)
		}
	}()

	router.GET("/", extendedScheduler)

	return router, server
}

func main() {
	fmt.Printf("Starting extended scheduler...")
	glog.V(4).Infof("Starting extended scheduler...")

	trustedPrefixConf := flag.String("trustedprefixconf","","config for scheduler")
	flag.Parse()
	fmt.Println("confpath is",*trustedPrefixConf)
	api.Confpath = *trustedPrefixConf
	fmt.Println("confpath is",api.Confpath)
	router, server := SetupRouter()

	//hadler for the post operation
	router.POST("filter", api.FilterHandler)

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	glog.V(4).Infof("Shutting down Extended Scheduler Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {

		glog.V(4).Infof("Extended Scheduler Server Shutdown:", err)
	}
	glog.V(4).Infof("Extended Scheduler Server exist")
	fmt.Printf("Stoping extended scheduler...")
}
