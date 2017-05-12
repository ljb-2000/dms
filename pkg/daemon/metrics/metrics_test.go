package metrics_test

import (
	"github.com/lavrs/dms/pkg/context"
	"github.com/lavrs/dms/pkg/daemon/docker"
	m "github.com/lavrs/dms/pkg/daemon/metrics"
	"github.com/lavrs/dms/pkg/logger"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	cName         = "splines"
	cImage        = "bfirsh/reticulate-splines"
	cAll          = "all"
	ucListTime    = time.Second * 1
	ucMetricsTime = time.Second * 1
)

func init() {
	err := docker.ImagePull(cImage)
	if err != nil {
		logger.Panic(err)
	}
	err = docker.ContainerCreate(cImage, cName)
	if err != nil {
		logger.Panic(err)
	}
	err = docker.ContainerStart(cName)
	if err != nil {
		logger.Panic(err)
	}

	context.Get().Debug = true

	m.Get().SetUCListInterval(ucListTime)
	m.Get().SetUCMetricsInterval(ucMetricsTime)
	go m.Get().Collect()
}

func pending(t time.Duration) {
	time.Sleep(t * 2)
}

func TestGetMetricsObj(t *testing.T) {
	assert.NotNil(t, m.Get())
}

func TestMetricsAlreadyCollecting(t *testing.T) {
	m.Get().Collect()
}

func TestGetNoRunningContainers(t *testing.T) {
	cMetrics := m.Get().Get(cName)
	assert.Equal(t, "no running containers", cMetrics.Message)
}

func TestGetContainerLogs(t *testing.T) {
	pending(ucListTime)

	logs := m.GetContainerLogs(cName)
	assert.NotEmpty(t, logs)
}

func TestGetLaunchedContainers(t *testing.T) {
	isChange := false

	launched := m.Get().GetLaunchedContainers()
	for i := range launched {
		if launched[i] == cName {
			isChange = true
		}
	}
	assert.True(t, isChange)
}

func TestGet(t *testing.T) {
	metrics := m.Get().Get(cAll)

	isChange := false
	for i := range metrics.Launched {
		if metrics.Launched[i] == cName {
			isChange = true
		}
	}
	assert.True(t, isChange)
}

func TestEmptyChanges(t *testing.T) {
	isChange := false
	launched := m.Get().GetLaunchedContainers()
	for i := range launched {
		if launched[i] == cName {
			isChange = true
		}
	}
	assert.False(t, isChange)

	isChange = false
	stopped := m.Get().GetStoppedContainers()
	for i := range stopped {
		if stopped[i] == cName {
			isChange = true
		}
	}
	assert.False(t, isChange)
}

func TestGetSpecifiedContainers(t *testing.T) {
	assert.Equal(t, "these containers are not running",
		m.Get().Get("container1 container2").Message)
}

func TestGetStoppedContainers(t *testing.T) {
	err := docker.ContainerStop(cName)
	assert.NoError(t, err)
	pending(ucListTime)

	isChange := false
	stopped := m.Get().GetStoppedContainers()
	for i := range stopped {
		if stopped[i] == cName {
			isChange = true
		}
	}
	assert.True(t, isChange)

	assert.Equal(t, cName, m.Get().Get(cName).Stopped[0])
}

func TestContainerRemoveHandle(t *testing.T) {
	err := docker.ContainerStart(cName)
	assert.NoError(t, err)
	pending(ucListTime)
	err = docker.ContainerRemove(cName)
	assert.NoError(t, err)
	pending(ucMetricsTime)
	err = docker.ImageRemove(cImage)
	assert.NoError(t, err)
}
