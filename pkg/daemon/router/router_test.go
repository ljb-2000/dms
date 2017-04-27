package router_test

import (
	"github.com/lavrs/docker-monitoring-service/pkg/context"
	"github.com/lavrs/docker-monitoring-service/pkg/daemon/router"
	"github.com/lavrs/docker-monitoring-service/pkg/docker"
	"github.com/stretchr/testify/assert"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/httptest"
	"testing"
)

func TestRun(t *testing.T) {
	const (
		cName  = "splines"
		cImage = "bfirsh/reticulate-splines"
	)

	context.Get().Debug = true

	e := httptest.New(router.App(), t)

	e.GET("/api/metrics/container").Expect().
		Status(iris.StatusOK).
		JSON().Object().Value("message").
		Equal("no running containers")

	e.GET("/api/logs/container").Expect().
		Status(iris.StatusOK).
		JSON().Object().Value("message").
		Equal("Error response from daemon: No such container: container")

	err := docker.ImagePull(cImage)
	assert.NoError(t, err)
	err = docker.ContainerCreate(cImage, cName)
	assert.NoError(t, err)
	err = docker.ContainerStart(cName)
	assert.NoError(t, err)

	e.GET("/api/logs/" + cName).Expect().
		Status(iris.StatusOK)

	err = docker.ContainerStop(cName)
	assert.NoError(t, err)
	err = docker.ContainerRemove(cName)
	assert.NoError(t, err)
	err = docker.ImageRemove(cImage)
	assert.NoError(t, err)

	e.GET("/charts").Expect().Header("Content-Type").
		Equal("text/html; charset=UTF-8")

	e.GET("/").Expect().Header("Content-Type").
		Equal("text/html; charset=UTF-8")
}
