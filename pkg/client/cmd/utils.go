package cmd

import (
	"encoding/json"
	"github.com/docker/docker/cli/command/formatter"
	h "github.com/lavrs/dms/pkg/client/http"
	"github.com/lavrs/dms/pkg/logger"
)

// GetContainersLogs returns container logs
func GetContainersLogs(url, id string) (string, error) {
	logger.Info("containers logs cmd")

	body, err := h.HTTPGET(url + "/logs/" + id)
	if err != nil {
		return "", err
	}

	var logs struct{ Logs string }
	err = json.Unmarshal(body, &logs)
	if err != nil {
		return "", err
	}

	return logs.Logs, nil
}

// GetContainersMetrics returns container(s) metrics
func GetContainersMetrics(url, id string) ([]*formatter.ContainerStats, error) {
	logger.Info("containers metrics cmd")

	body, err := h.HTTPGET(url + "/metrics/" + id)
	if err != nil {
		return nil, err
	}

	var metrics struct{ Metrics []*formatter.ContainerStats }
	err = json.Unmarshal(body, &metrics)
	if err != nil {
		return nil, err
	}

	return metrics.Metrics, nil
}

// GetStoppedContainers returns stopped containers
func GetStoppedContainers(url string) ([]string, error) {
	logger.Info("stopped containers cmd")

	body, err := h.HTTPGET(url + "/stopped")
	if err != nil {
		return nil, err
	}

	var stopped struct{ Stopped []string }
	err = json.Unmarshal(body, &stopped)
	if err != nil {
		return nil, err
	}

	return stopped.Stopped, nil
}

// GetLaunchedContainers returns launched containers
func GetLaunchedContainers(url string) ([]string, error) {
	logger.Info("launched containers cmd")

	body, err := h.HTTPGET(url + "/launched")
	if err != nil {
		return nil, err
	}

	var launched struct{ Launched []string }
	err = json.Unmarshal(body, &launched)
	if err != nil {
		return nil, err
	}

	return launched.Launched, nil
}
