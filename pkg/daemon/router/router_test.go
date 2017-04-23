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
		JSON().Object().Value("message").Equal("no running containers")
}
