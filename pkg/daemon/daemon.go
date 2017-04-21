package daemon

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	m "github.com/lavrs/docker-monitoring-service/pkg/metrics"
	"html/template"
	"io"
	"net/http"
	"time"
)

func Run(port string, ucltime, uctime int) error {
	const (
		rootDir       = "website"
		indexTemplate = "index.html"
	)

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.Static(rootDir))

	metrics := m.NewMetrics()
	metrics.SetUCLTime(time.Duration(ucltime) * time.Second)
	metrics.SetUCTime(time.Duration(uctime) * time.Second)

	go metrics.Collect()

	e.GET("/metrics/:id", func(c echo.Context) error {
		return c.JSON(http.StatusOK, metrics.Get(c.Param("id")))
	})

	e.GET("/charts", func(c echo.Context) error {
		t := &Template{
			templates: template.Must(template.ParseGlob(rootDir + "/" + indexTemplate)),
		}
		e.Renderer = t

		return c.Render(http.StatusOK, indexTemplate, nil)
	})

	return e.Start(":" + port)
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
