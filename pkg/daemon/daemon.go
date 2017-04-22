package daemon

import (
	dmsLogger "github.com/lavrs/docker-monitoring-service/pkg/logger"
	m "github.com/lavrs/docker-monitoring-service/pkg/metrics"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/cors"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/view"
	iLogger "gopkg.in/kataras/iris.v6/middleware/logger"
	"gopkg.in/kataras/iris.v6/middleware/recover"
	"net/http"
	"time"
)

// run daemon
func Run(port string, uctl, ucl int) error {
	metrics, err := m.NewMetrics()
	if err != nil {
		dmsLogger.Panic(err)
	}
	metrics.SetUCLTime(time.Duration(uctl) * time.Second)
	metrics.SetUCTime(time.Duration(ucl) * time.Second)
	go metrics.Collect()

	fsrv := &http.Server{
		Handler: App(),
		Addr:    ":" + port,
	}
	return fsrv.ListenAndServe()
}

// get daemon configuration
func App() *iris.Framework {
	return app()
}

func app() *iris.Framework {
	app := iris.New()
	app.Adapt(
		httprouter.New(),
		iris.DevLogger(),
		cors.New(cors.Options{AllowedOrigins: []string{"*"}}),
		view.HTML("./website", ".html"),
	)
	app.Use(
		recover.New(),
		iLogger.New(iLogger.Config{
			Status: true,
			IP:     true,
			Method: true,
			Path:   true,
		}),
	)
	app.StaticWeb("/static", "website/static")

	app.Get("/api/metrics/:id", getMetrics)
	app.OnError(iris.StatusNotFound, p404)
	app.Get("/charts", charts)

	app.Boot()
	return app
}

func charts(ctx *iris.Context) {
	ctx.MustRender("index.html", nil)
}

func p404(ctx *iris.Context) {
	ctx.MustRender("404.html", nil)
}

func getMetrics(ctx *iris.Context) {
	ctx.JSON(iris.StatusOK, m.Get().Get(ctx.Param("id")))
}
