package docker

import (
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/lavrs/dms/pkg/daemon/docker"
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
	reader, err := docker.ContainerStats(cName)
	assert.NoError(t, err)
	assert.NotNil(t, reader)
}

func TestFormatting(t *testing.T) {
	reader, err := docker.ContainerStats(cName)
	assert.NoError(t, err)

	dec := json.NewDecoder(reader)
	var statsJSON *types.StatsJSON
	err = dec.Decode(&statsJSON)
	assert.NoError(t, err)

	assert.NotNil(t, docker.Formatting(statsJSON))
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
