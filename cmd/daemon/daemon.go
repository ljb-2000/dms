package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	dms "github.com/lavrs/docker-monitoring-service/stats"
	"net/http"
)

func main() {
	e := echo.New()

	var s dms.Stats

	go s.CollectData()

	e.GET("/stats/:id", func(c echo.Context) error {
		stats, running, stopped, err := s.Get(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"stats":              stats,
			"running_containers": running,
			"stopped_containers": stopped,
		})
	})

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Logger.Fatal(e.Start(":8080"))
}
