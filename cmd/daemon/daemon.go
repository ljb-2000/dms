package main

import (
	"context"
	"github.com/docker/docker/client"
	s "github.com/lavrs/docker-monitoring-service/stats"
	"gopkg.in/gin-gonic/gin.v1"
)

func main() {
	r := gin.Default()

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	r.GET("/", func(c *gin.Context) {
		stats := s.Stats(context.Background(), cli, "splines")

		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, gin.H{
			"cpu": stats.CPUPercentage,
			"mem": stats.MemoryPercentage,
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run()
}
