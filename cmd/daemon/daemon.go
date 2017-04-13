package main

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	s "github.com/lavrs/docker-monitoring-service/stats"
	"net/http"
)

func main() {
	e := echo.New()

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "OK",
		})
	})

	e.GET("/:id", func(c echo.Context) error {
		stats := s.Stats(context.Background(), cli, c.Param("id"))

		return c.JSON(http.StatusOK, map[string]interface{}{
			"cpu": stats.GetStatistics().CPUPercentage,
			"mem": stats.GetStatistics().MemoryPercentage,
		})
	})

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Logger.Fatal(e.Start(":8080"))
}
