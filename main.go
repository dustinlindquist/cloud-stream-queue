package main

import (
	"cloud-stream-queue/api"
	"cloud-stream-queue/queue"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1"})

	queue := queue.New()
	api := api.New(queue, router)
	api.Router.Run(":3000")
}
