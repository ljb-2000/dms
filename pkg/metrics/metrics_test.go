package metrics_test

import (
	"github.com/lavrs/docker-monitoring-service/pkg/context"
	"github.com/lavrs/docker-monitoring-service/pkg/docker"
	m "github.com/lavrs/docker-monitoring-service/pkg/metrics"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewMetrics(t *testing.T) {
	metrics := m.Get()
	assert.NotNil(t, metrics)
}

func TestGet_Collect(t *testing.T) {
	const (
		cName         = "splines"
		cImage        = "bfirsh/reticulate-splines"
		cAll          = "all"
		ucListTime    = time.Second * 1
		ucMetricsTime = time.Second * 1
	)

	var (
		isChange = false
		ctx      = context.Get()
	)

	ctx.Debug = true

	metrics := m.Get()
	metrics.SetUCListInterval(ucListTime)
	metrics.SetUCMetricsInterval(ucMetricsTime)
	go metrics.Collect()

	cMetrics := metrics.Get(cName)
	assert.Equal(t, "no running containers", cMetrics.Message)

	err := docker.ImagePull(cImage)
	assert.NoError(t, err)
	err = docker.ContainerCreate(cImage, cName)
	assert.NoError(t, err)
	err = docker.ContainerStart(cName)
	assert.NoError(t, err)
	pending(ucListTime)

	logs := m.GetContainerLogs(cName)
	assert.NotEmpty(t, logs)

    isChange = false
    launched := metrics.GetLaunchedContainers()
    for i := range launched {
        if launched[i] == cName {
            isChange = true
        }
    }
    assert.True(t, isChange)

	cMetrics = metrics.Get(cAll)
	for i := range cMetrics.Launched {
		if cMetrics.Launched[i] == cName {
			isChange = true
		}
	}
	assert.True(t, isChange)

    isChange = false
    launched = metrics.GetLaunchedContainers()
    for i := range launched {
        if launched[i] == cName {
            isChange = true
        }
    }
    assert.False(t, isChange)

	isChange = false
	stopped := metrics.GetStoppedContainers()
	for i := range stopped {
		if stopped[i] == cName {
			isChange = true
		}
	}
	assert.False(t, isChange)

	err = docker.ContainerStop(cName)
	assert.NoError(t, err)
	err = docker.ContainerStart(cName)
	assert.NoError(t, err)
	pending(ucListTime)

	cMetrics = metrics.Get("container1 container2")
	assert.Equal(t, "these containers are not running", cMetrics.Message)

	err = docker.ContainerStop(cName)
	assert.NoError(t, err)
	pending(ucListTime)

	isChange = false
	stopped = metrics.GetStoppedContainers()
	for i := range stopped {
		if stopped[i] == cName {
			isChange = true
		}
	}
	assert.True(t, isChange)

	cMetrics = metrics.Get(cAll)
	assert.Equal(t, cName, cMetrics.Stopped[0])

	err = docker.ContainerStart(cName)
	assert.NoError(t, err)
	pending(ucListTime)
	err = docker.ContainerRemove(cName)
	assert.NoError(t, err)
	pending(ucListTime)
	err = docker.ImageRemove(cImage)
	assert.NoError(t, err)
}

func pending(t time.Duration) {
	time.Sleep(t * 2)
}
