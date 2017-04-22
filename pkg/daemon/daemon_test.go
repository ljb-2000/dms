package daemon_test

import (
	"github.com/lavrs/docker-monitoring-service/pkg/daemon"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/httptest"
	"testing"
)

func TestRun(t *testing.T) {
	e := httptest.New(daemon.App(), t)

	e.GET("/api/metrics/container").Expect().
		Status(iris.StatusOK).
		JSON().Object().Value("message").Equal("no running containers")
}
