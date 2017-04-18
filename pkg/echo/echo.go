package echo

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/lavrs/docker-monitoring-service/pkg/metrics"
	"net/http"
)

func Echo(port string) {
	e := echo.New()

	m := metrics.NewMetrics()

	go m.Collect()

	e.GET("/metrics/:id", func(c echo.Context) error {
		return c.JSON(http.StatusOK, m.Get(c.Param("id")))
	})

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Logger.Fatal(e.Start(":" + port))
}
