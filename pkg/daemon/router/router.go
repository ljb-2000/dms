package router

import (
	"context"
	ctx "github.com/lavrs/dms/pkg/context"
	"github.com/lavrs/dms/pkg/daemon/metrics"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/cors"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/view"
	"gopkg.in/kataras/iris.v6/middleware/logger"
	"gopkg.in/kataras/iris.v6/middleware/recover"
	"time"
)

// App returns daemon configuration
func App() *iris.Framework {
	return app()
}

// iris configuration
func app() *iris.Framework {
	app := iris.New()
	app.Adapt(
		iris.EventPolicy{
			Interrupted: func(*iris.Framework) {
				ctxwt, _ := context.WithTimeout(context.Background(), 1*time.Second)
				app.Shutdown(ctxwt)
			},
		},
		httprouter.New(),
		cors.New(cors.Options{AllowedOrigins: []string{"*"}}),
		view.HTML("./dashboard", ".html"),
	)
	app.StaticWeb("/static", "dashboard/static")
	if ctx.Get().Debug {
		app.Use(
			recover.New(),
			logger.New(logger.Config{
				Status: true,
				IP:     true,
				Method: true,
				Path:   true,
			}),
		)
		app.Adapt(iris.DevLogger())
	}

	app.Get("/dashboard", charts)
	app.OnError(iris.StatusNotFound, p404)

	app.Get("/status", status)

	app.Get("/api/logs/:id", getLogs)
	app.Get("/api/metrics/:id", getMetrics)
	app.Get("/api/stopped", getStopped)
	app.Get("/api/launched", getLaunched)

	app.Boot()
	return app
}

func status(ctx *iris.Context) {
	ctx.WriteHeader(iris.StatusOK)
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

// get stopped containers
func getStopped(ctx *iris.Context) {
	ctx.JSON(iris.StatusOK, map[string][]string{
		"stopped": metrics.Get().GetStoppedContainers(),
	})
}

// get launched containers
func getLaunched(ctx *iris.Context) {
	ctx.JSON(iris.StatusOK, map[string][]string{
		"launched": metrics.Get().GetLaunchedContainers(),
	})
}

// get container logs
func getLogs(ctx *iris.Context) {
	ctx.JSON(iris.StatusOK, map[string]string{
		"logs": metrics.GetContainerLogs(ctx.Param("id")),
	})
}
