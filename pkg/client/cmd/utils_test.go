package cmd_test

import (
	"encoding/json"
	"github.com/docker/docker/cli/command/formatter"
	"github.com/lavrs/dms/pkg/client/cmd"
	"github.com/lavrs/dms/pkg/context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetContainerLogs(t *testing.T) {
	const testLogs = "test logs"

	var isInternalServerError = false

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logs, err := json.Marshal(map[string]string{
			"logs": testLogs,
		})
		assert.NoError(t, err)

		if !isInternalServerError {
			isInternalServerError = true

			w.WriteHeader(http.StatusOK)
			_, err = w.Write(logs)
			assert.NoError(t, err)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer ts.Close()

	context.Get().Address = ts.URL

	logs, err := cmd.GetContainersLogs("container")
	assert.NoError(t, err)
	assert.Equal(t, testLogs, logs)

	_, err = cmd.GetContainersLogs("container")
	assert.Error(t, err)
}

func TestGetStoppedContainers(t *testing.T) {
	var (
		testStopped           = []string{"container1", "container2"}
		isInternalServerError = false
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stopped, err := json.Marshal(map[string][]string{
			"stopped": testStopped,
		})
		assert.NoError(t, err)

		if !isInternalServerError {
			isInternalServerError = true

			w.WriteHeader(http.StatusOK)
			_, err = w.Write(stopped)
			assert.NoError(t, err)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer ts.Close()

	context.Get().Address = ts.URL

	stopped, err := cmd.GetStoppedContainers()
	assert.NoError(t, err)
	assert.Equal(t, testStopped, stopped)

	_, err = cmd.GetStoppedContainers()
	assert.Error(t, err)
}

func TestGetLaunchedContainers(t *testing.T) {
	var (
		testLaunched          = []string{"container1", "container2"}
		isInternalServerError = false
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		launched, err := json.Marshal(map[string][]string{
			"launched": testLaunched,
		})
		assert.NoError(t, err)

		if !isInternalServerError {
			isInternalServerError = true

			w.WriteHeader(http.StatusOK)
			_, err = w.Write(launched)
			assert.NoError(t, err)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer ts.Close()

	context.Get().Address = ts.URL

	launched, err := cmd.GetLaunchedContainers()
	assert.NoError(t, err)
	assert.Equal(t, testLaunched, launched)

	_, err = cmd.GetLaunchedContainers()
	assert.Error(t, err)
}

func TestGetContainersMetrics(t *testing.T) {
	const (
		container1 = "container1"
		container2 = "container2"
	)

	var isInternalServerError = false

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

		if !isInternalServerError {
			isInternalServerError = true

			w.WriteHeader(http.StatusOK)
			_, err = w.Write(metrics)
			assert.NoError(t, err)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer ts.Close()

	context.Get().Address = ts.URL

	metrics, err := cmd.GetContainersMetrics("all")
	assert.NoError(t, err)
	assert.Equal(t, container1, metrics[0].ID)
	assert.Equal(t, container2, metrics[1].ID)

	_, err = cmd.GetContainersMetrics("all")
	assert.Error(t, err)
}
