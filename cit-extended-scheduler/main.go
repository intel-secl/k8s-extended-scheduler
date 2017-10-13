package main

import (
	"fmt"
	"cit-extended-scheduler/api"
	"github.com/gin-gonic/gin"
)

// PORT is the port to listen.
const PORT = 8888

func main() {
	r := gin.Default()
	fmt.Printf("Starting extended scheduler...")
	r.POST("filter", api.FilterHandler)
	r.Run(fmt.Sprintf(":%d", PORT))
}
