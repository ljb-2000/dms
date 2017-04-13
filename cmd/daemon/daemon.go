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

	e.GET("/get/all", func(c echo.Context) error {
		return c.JSON(http.StatusOK, s.AllStats(cli))
	})

	e.GET("/get/:id", func(c echo.Context) error {
		return c.JSON(http.StatusOK, s.Stats(cli, c.Param("id")))
	})

	e.Any("*", func(c echo.Context) error {
		return c.String(http.StatusNotFound, "page not found")
	})

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Logger.Fatal(e.Start(":8080"))
}
