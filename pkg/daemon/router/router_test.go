package router_test

import (
	"github.com/lavrs/docker-monitoring-service/pkg/context"
	"github.com/lavrs/docker-monitoring-service/pkg/daemon/router"
	"github.com/stretchr/testify/assert"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/httptest"
	"testing"
)

var e = httptest.New(router.App(), nil)

func TestApp(t *testing.T) {
	context.Get().Debug = true

	e = httptest.New(router.App(), t)
	assert.NotNil(t, e)
}

func TestApiMetrics(t *testing.T) {
	e.GET("/api/metrics/container").Expect().
		Status(iris.StatusOK).
		JSON().Object().Value("message").
		Equal("no running containers")
}

func TestApiLogs(t *testing.T) {
	e.GET("/api/logs/container").Expect().Status(iris.StatusInternalServerError)
}

func TestApiStopped(t *testing.T) {
	e.GET("/api/stopped").Expect().
		Status(iris.StatusOK)
}

func TestApiLaunched(t *testing.T) {
	e.GET("/api/launched").Expect().
		Status(iris.StatusOK)
}

func TestChartsPage(t *testing.T) {
	e.GET("/charts").Expect().Header("Content-Type").
		Equal("text/html; charset=UTF-8")
}

func Test404Page(t *testing.T) {
	e.GET("/").Expect().Header("Content-Type").
		Equal("text/html; charset=UTF-8")
}
