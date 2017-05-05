package cmd_test

import (
	"encoding/json"
	"github.com/docker/docker/cli/command/formatter"
	"github.com/lavrs/dms/pkg/client/cmd"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetContainerLogs(t *testing.T) {
	const testLogs = "test logs"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logs, err := json.Marshal(map[string]string{
			"logs": testLogs,
		})
		assert.NoError(t, err)

		w.WriteHeader(200)
		_, err = w.Write(logs)
		assert.NoError(t, err)
	}))
	defer ts.Close()

	logs, err := cmd.GetContainersLogs(ts.URL, "container")
	assert.NoError(t, err)
	assert.Equal(t, testLogs, logs)
}

func TestGetStoppedContainers(t *testing.T) {
	var testStopped = []string{"container1", "container2"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stopped, err := json.Marshal(map[string][]string{
			"stopped": testStopped,
		})
		assert.NoError(t, err)

		w.WriteHeader(200)
		_, err = w.Write(stopped)
		assert.NoError(t, err)
	}))
	defer ts.Close()

	stopped, err := cmd.GetStoppedContainers(ts.URL)
	assert.NoError(t, err)
	assert.Equal(t, testStopped, stopped)
}

func TestGetLaunchedContainers(t *testing.T) {
	var testLaunched = []string{"container1", "container2"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		launched, err := json.Marshal(map[string][]string{
			"launched": testLaunched,
		})
		assert.NoError(t, err)

		w.WriteHeader(200)
		_, err = w.Write(launched)
		assert.NoError(t, err)
	}))
	defer ts.Close()

	launched, err := cmd.GetLaunchedContainers(ts.URL)
	assert.NoError(t, err)
	assert.Equal(t, testLaunched, launched)
}

func TestGetContainersMetrics(t *testing.T) {
	const (
		container1 = "container1"
		container2 = "container2"
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		w.WriteHeader(200)
		_, err = w.Write(metrics)
		assert.NoError(t, err)
	}))
	defer ts.Close()

	metrics, err := cmd.GetContainersMetrics(ts.URL, "all")
	assert.NoError(t, err)
	assert.Equal(t, container1, metrics[0].ID)
	assert.Equal(t, container2, metrics[1].ID)
}
