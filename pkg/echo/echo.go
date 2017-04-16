package echo

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/lavrs/docker-monitoring-service/pkg/stats"
	"net/http"
)

func Echo(port string) {
	e := echo.New()

	var s stats.Stats
	go s.Collect()

	e.GET("/stats/:id", func(c echo.Context) error {
		return c.JSON(http.StatusOK, s.Get(c.Param("id")))
	})

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Logger.Fatal(e.Start(":" + port))
}
