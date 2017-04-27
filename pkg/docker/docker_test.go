package docker_test

import (
	"github.com/docker/docker/cli/command/formatter"
	"github.com/lavrs/docker-monitoring-service/pkg/docker"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	cName  = "splines"
	cImage = "bfirsh/reticulate-splines"
)

func TestImagePull(t *testing.T) {
	err := docker.ImagePull(cImage)
	assert.NoError(t, err)
}

func TestContainerCreate(t *testing.T) {
	err := docker.ContainerCreate(cImage, cName)
	assert.NoError(t, err)
}

func TestContainerStart(t *testing.T) {
	err := docker.ContainerStart(cName)
	assert.NoError(t, err)
}

func TestContainerList(t *testing.T) {
	container, err := docker.ContainerList()
	assert.NoError(t, err)
	assert.NotNil(t, container)
}

func TestContainerStats(t *testing.T) {
	var (
		err error

		data  = make(chan *formatter.ContainerStats)
		gDone = make(chan bool)
		done  = make(chan bool)
	)

	go func() {
		err = docker.ContainerStats(cName, data, done)
		close(gDone)
	}()

	<-data
	close(done)
	close(data)
	<-gDone

	assert.NoError(t, err)
}

func TestContainersLogs(t *testing.T) {
	_, err := docker.ContainersLogs(cName)
	assert.NoError(t, err)
}

func TestContainerStop(t *testing.T) {
	err := docker.ContainerStop(cName)
	assert.NoError(t, err)
}

func TestContainerRemove(t *testing.T) {
	err := docker.ContainerRemove(cName)
	assert.NoError(t, err)
}

func TestImageRemove(t *testing.T) {
	err := docker.ImageRemove(cImage)
	assert.NoError(t, err)
}
