package metrics_test

import (
	"github.com/lavrs/docker-monitoring-service/pkg/docker"
	m "github.com/lavrs/docker-monitoring-service/pkg/metrics"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewMetrics(t *testing.T) {
	metrics := m.NewMetrics()
	assert.NotNil(t, metrics)
}

func TestGet_Collect(t *testing.T) {
	const (
		cName         = "splines"
		dockerCName   = "/splines"
		cImage        = "bfirsh/reticulate-splines"
		cAll          = "all"
		ucListTime    = time.Second * 2
		ucMetricsTime = time.Second * 1
	)

	var (
		isLaunched = false
	)

	metrics := m.NewMetrics()
	assert.NotNil(t, metrics)
	metrics.SetUCLTime(ucListTime)
	metrics.SetUCTime(ucMetricsTime)

	go metrics.Collect()

	cMetrics := metrics.Get(cName)
	assert.Equal(t, "no running containers", cMetrics.Message)

	err := docker.ImagePull(cImage)
	assert.NoError(t, err)
	err = docker.ContainerCreate(cImage, cName)
	assert.NoError(t, err)
	err = docker.ContainerStart(cName)
	assert.NoError(t, err)
	defer func() {
		err = docker.ContainerRemove(cName)
		assert.NoError(t, err)

		err = docker.ImageRemove(cImage)
		assert.NoError(t, err)
	}()
	pending(ucListTime)

	cMetrics = metrics.Get(cAll)
	for i, _ := range cMetrics.Launched {
		if cMetrics.Launched[i] == dockerCName {
			isLaunched = true
		}
	}
	assert.True(t, isLaunched)

	cMetrics = metrics.Get("container1 container2")
	assert.Equal(t, "these containers are not running", cMetrics.Message)

	err = docker.ContainerStop(cName)
	assert.NoError(t, err)
	pending(ucListTime)

	cMetrics = metrics.Get(cName)
	assert.Equal(t, dockerCName, cMetrics.Stopped[0])
}

func pending(t time.Duration) {
	time.Sleep(t * 2)
}
