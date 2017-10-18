package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
	"cit-extended-scheduler/api"
	"cit-extended-scheduler/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)


func main() {
	fmt.Printf("Starting extended scheduler...")
	glog.V(4).Infof("Starting extended scheduler...")

	//get a webserver instance, that contains a muxer, middleware and configuration settings
	router := gin.Default()

	// fetch all the cmd line args
	url, port, server_crt, server_key := util.GetCmdlineArgs()

	//initialize http server config
	srv := &http.Server{
                Addr:    *url+":"+port,
                Handler: router,
        }

	//hadler for the post operation 
	router.POST("filter", api.FilterHandler)

	//run the server instance
	go func() {
                // service connections
                if err := srv.ListenAndServeTLS(*server_crt, *server_key); err != nil {
                        glog.V(4).Infof("listen: %s\n", err)
			fmt.Printf("listen %s ...", err)
                }
        }()
	
	// Wait for interrupt signal to gracefully shutdown the server with
        // a timeout of 5 seconds.
        quit := make(chan os.Signal)
        signal.Notify(quit, os.Interrupt)
        <-quit
        glog.V(4).Infof("Shutting down Extended Scheduler Server ...")

        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        if err := srv.Shutdown(ctx); err != nil {

                glog.V(4).Infof("Extended Scheduler Server Shutdown:", err)
        }
        glog.V(4).Infof("Extended Scheduler Server exist")
	fmt.Printf("Stoping extended scheduler...")
}
