package cmd_test

import (
	"encoding/json"
	"github.com/docker/cli/cli/command/formatter"
	"github.com/lavrs/dms/pkg/client/cmd"
	"github.com/lavrs/dms/pkg/context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetContainerLogs(t *testing.T) {
	const (
		testLogs    = "test logs"
		containerID = "container"
	)

	ts200 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logs, err := json.Marshal(map[string]string{
			"logs": testLogs,
		})
		assert.NoError(t, err)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(logs)
		assert.NoError(t, err)
	}))
	defer ts200.Close()

	context.Get().Address = ts200.URL
	cmd.ContainersLogsCmd(containerID)
	logs, err := cmd.GetContainersLogs(containerID)
	assert.NoError(t, err)
	assert.Equal(t, testLogs, logs)

	ts500 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts500.Close()

	context.Get().Address = ts500.URL
	cmd.ContainersLogsCmd(containerID)
	_, err = cmd.GetContainersLogs(containerID)
	assert.Error(t, err)
}

func TestGetStoppedContainers(t *testing.T) {
	var (
		testStopped         = []string{"container1", "container2"}
		noStoppedContainers = []string{"no stopped containers"}
	)

	ts200 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stopped, err := json.Marshal(map[string][]string{
			"stopped": testStopped,
		})
		assert.NoError(t, err)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(stopped)
		assert.NoError(t, err)
	}))
	defer ts200.Close()

	context.Get().Address = ts200.URL
	cmd.StoppedContainersCmd()
	stopped, err := cmd.GetStoppedContainers()
	assert.NoError(t, err)
	assert.Equal(t, testStopped, stopped)

	ts200NoStoppedContainers := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stopped, err := json.Marshal(map[string][]string{
			"stopped": noStoppedContainers,
		})
		assert.NoError(t, err)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(stopped)
		assert.NoError(t, err)
	}))
	defer ts200NoStoppedContainers.Close()

	context.Get().Address = ts200NoStoppedContainers.URL
	cmd.StoppedContainersCmd()

	ts500 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts500.Close()

	context.Get().Address = ts500.URL
	cmd.StoppedContainersCmd()
	_, err = cmd.GetStoppedContainers()
	assert.Error(t, err)
}

func TestGetLaunchedContainers(t *testing.T) {
	var (
		testLaunched         = []string{"container1", "container2"}
		noLaunchedContainers = []string{"no launched containers"}
	)

	ts200 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		launched, err := json.Marshal(map[string][]string{
			"launched": testLaunched,
		})
		assert.NoError(t, err)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(launched)
		assert.NoError(t, err)
	}))
	defer ts200.Close()

	context.Get().Address = ts200.URL
	cmd.LaunchedContainersCmd()
	launched, err := cmd.GetLaunchedContainers()
	assert.NoError(t, err)
	assert.Equal(t, testLaunched, launched)

	ts200NoLaunchedContainers := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stopped, err := json.Marshal(map[string][]string{
			"launched": noLaunchedContainers,
		})
		assert.NoError(t, err)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(stopped)
		assert.NoError(t, err)
	}))
	defer ts200NoLaunchedContainers.Close()

	context.Get().Address = ts200NoLaunchedContainers.URL
	cmd.LaunchedContainersCmd()

	ts500 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts500.Close()

	context.Get().Address = ts500.URL
	cmd.LaunchedContainersCmd()
	_, err = cmd.GetLaunchedContainers()
	assert.Error(t, err)
}

func TestGetContainersMetrics(t *testing.T) {
	const (
		container1 = "container1"
		container2 = "container2"
	)

	var all = []string{"all"}

	ts200 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			cS1   formatter.ContainerStats
			cS2   formatter.ContainerStats
			stats []*formatter.ContainerStats
		)

		cS1.SetStatistics(formatter.StatsEntry{
			ID: container1,
		})
		cS2.SetStatistics(formatter.StatsEntry{
			ID: container2,
		})
		stats = append(stats, &cS1)
		stats = append(stats, &cS2)

		metrics, err := json.Marshal(map[string]*[]*formatter.ContainerStats{
			"metrics": &stats,
		})
		assert.NoError(t, err)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(metrics)
		assert.NoError(t, err)
	}))
	defer ts200.Close()

	context.Get().Address = ts200.URL
	cmd.ContainersMetricsCmd(all)
	metrics, err := cmd.GetContainersMetrics("all")
	assert.NoError(t, err)
	assert.Equal(t, container1, metrics[0].ID)
	assert.Equal(t, container2, metrics[1].ID)

	ts200NilMetrics := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var stats []*formatter.ContainerStats

		metrics, err := json.Marshal(map[string]*[]*formatter.ContainerStats{
			"metrics": &stats,
		})
		assert.NoError(t, err)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(metrics)
		assert.NoError(t, err)
	}))
	defer ts200NilMetrics.Close()

	context.Get().Address = ts200NilMetrics.URL
	cmd.ContainersMetricsCmd(all)

	ts500 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts500.Close()

	context.Get().Address = ts500.URL
	cmd.ContainersMetricsCmd(all)
	_, err = cmd.GetContainersMetrics("all")
	assert.Error(t, err)
}
