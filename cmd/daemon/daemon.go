package main

import (
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

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "pong",
		})
	})

	e.GET("/stats/:id", func(c echo.Context) error {
		stats, err := s.GetStats(cli, c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusOK, stats)
	})

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Logger.Fatal(e.Start(":8080"))
}
