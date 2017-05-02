package router_test

import (
	"github.com/lavrs/docker-monitoring-service/pkg/context"
	"github.com/lavrs/docker-monitoring-service/pkg/daemon/router"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/httptest"
	"testing"
)

func TestRun(t *testing.T) {
	context.Get().Debug = true

	e := httptest.New(router.App(), t)

	e.GET("/api/metrics/container").Expect().
		Status(iris.StatusOK).
		JSON().Object().Value("message").
		Equal("no running containers")

	e.GET("/api/logs/container").Expect().Status(iris.StatusInternalServerError)

	e.GET("/api/stopped").Expect().
		Status(iris.StatusOK)

	e.GET("/api/launched").Expect().
		Status(iris.StatusOK)

	e.GET("/charts").Expect().Header("Content-Type").
		Equal("text/html; charset=UTF-8")

	e.GET("/").Expect().Header("Content-Type").
		Equal("text/html; charset=UTF-8")
}
