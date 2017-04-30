package router

import (
	"github.com/lavrs/docker-monitoring-service/pkg/context"
	"github.com/lavrs/docker-monitoring-service/pkg/docker"
	"github.com/lavrs/docker-monitoring-service/pkg/metrics"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/cors"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/view"
	"gopkg.in/kataras/iris.v6/middleware/logger"
	"gopkg.in/kataras/iris.v6/middleware/recover"
)

// get daemon configuration
func App() *iris.Framework {
	return app()
}

// iris configuration
func app() *iris.Framework {
	app := iris.New()
	app.Adapt(
		httprouter.New(),
		cors.New(cors.Options{AllowedOrigins: []string{"*"}}),
		view.HTML("./website", ".html"),
	)
	app.Use(recover.New())
	app.StaticWeb("/static", "website/static")
	if context.Get().Debug {
		app.Use(
			logger.New(logger.Config{
				Status: true,
				IP:     true,
				Method: true,
				Path:   true,
			}),
		)
		app.Adapt(iris.DevLogger())
	}

	app.Get("/api/metrics/:id", getMetrics)
	app.OnError(iris.StatusNotFound, p404)
	app.Get("/charts", charts)
	app.Get("/api/logs/:id", getLogs)

	app.Boot()
	return app
}

// charts page
func charts(ctx *iris.Context) {
	ctx.MustRender("index.html", nil)
}

// 404 page
func p404(ctx *iris.Context) {
	ctx.MustRender("404.html", nil)
}

// get container metrics
func getMetrics(ctx *iris.Context) {
	ctx.JSON(iris.StatusOK, metrics.Get().Get(ctx.Param("id")))
}

// get container logs
func getLogs(ctx *iris.Context) {
	logs, err := docker.ContainersLogs(ctx.Param("id"))
	if err != nil {
		ctx.JSON(iris.StatusOK, map[string]string{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(iris.StatusOK, logs)
}
