package main

import (
	"github.com/gin-gonic/gin"
	// "log"
	"bus-api/src"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})
	r.GET("/bus-timing", src.GetBusTiming)
	r.GET("/list-of-bus-stop", src.GetListOfBusStop)
	r.Run()
}

