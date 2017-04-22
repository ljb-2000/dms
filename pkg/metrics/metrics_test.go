package metrics_test

import (
	"github.com/lavrs/docker-monitoring-service/pkg/context"
	"github.com/lavrs/docker-monitoring-service/pkg/docker"
	m "github.com/lavrs/docker-monitoring-service/pkg/metrics"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var metrics, err = m.NewMetrics()

func TestNewMetrics(t *testing.T) {
	assert.NoError(t, err)
	assert.NotNil(t, metrics)
}

func TestMultMetrics(t *testing.T) {
	_, err = m.NewMetrics()
	assert.Error(t, err)
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
		isLaunched = false
		ctx        = context.Get()
	)

	ctx.Debug = false

	metrics = m.Get()
	metrics.SetUCLTime(ucListTime)
	metrics.SetUCTime(ucMetricsTime)

	go metrics.Collect()

	cMetrics := metrics.Get(cName)
	assert.Equal(t, "no running containers", cMetrics.Message)

	err = docker.ImagePull(cImage)
	assert.NoError(t, err)
	err = docker.ContainerCreate(cImage, cName)
	assert.NoError(t, err)
	err = docker.ContainerStart(cName)
	assert.NoError(t, err)
	pending(ucListTime)

	cMetrics = metrics.Get(cAll)
	for i, _ := range cMetrics.Launched {
		if cMetrics.Launched[i] == cName {
			isLaunched = true
		}
	}
	assert.True(t, isLaunched)

	cMetrics = metrics.Get("container1 container2")
	assert.Equal(t, "these containers are not running", cMetrics.Message)

	err = docker.ContainerStop(cName)
	assert.NoError(t, err)
	pending(ucListTime)

	cMetrics = metrics.Get(cAll)
	assert.Equal(t, cName, cMetrics.Stopped[0])

	err = docker.ContainerStart(cName)
	assert.NoError(t, err)
	pending(ucListTime)
	err = docker.ContainerStop(cName)
	assert.NoError(t, err)
	err = docker.ContainerRemove(cName)
	assert.NoError(t, err)
	err = docker.ImageRemove(cImage)
	assert.NoError(t, err)
}

func pending(t time.Duration) {
	time.Sleep(t * 2)
}
