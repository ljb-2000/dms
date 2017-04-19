package metrics

import (
	"github.com/lavrs/docker-monitoring-service/pkg/docker"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	cName         = "splines"
	cImage        = "bfirsh/reticulate-splines"
	cAll          = "all"
	ucListTime    = time.Second * 3
	ucMetricsTime = time.Second * 1
)

func TestNewMetrics(t *testing.T) {
	m := NewMetrics()
	assert.NotNil(t, m)
}

func TestGet_Collect(t *testing.T) {
	m := NewMetrics()
	assert.NotNil(t, m)

	go m.Collect()

	metrics := m.Get(cName)
	assert.Equal(t, "no running containers", metrics.Message)

	err := docker.ImagePull(cImage)
	assert.NoError(t, err)

	err = docker.StartContainer(cImage, cName)
	assert.NoError(t, err)
	defer func() {
		err = docker.RemoveContainer(cName)
		assert.NoError(t, err)

		err = docker.ImageRemove(cImage)
		assert.NoError(t, err)
	}()
	pending()

	metrics = m.Get(cAll)
	assert.Equal(t, cName, metrics.Launched[0])
	assert.Equal(t, cName, metrics.Metrics[0].Name)

	metrics = m.Get("container1 container2")
	assert.Equal(t, "these containers are not running", metrics.Message)

	err = docker.StopContainer(cName)
	assert.NoError(t, err)
	pending()

	metrics = m.Get(cName)
	assert.Equal(t, cName, metrics.Stopped[0])
}

func pending() {
	time.Sleep(ucListTime * 2)
}
