package cmd

import (
	"encoding/json"
	"github.com/docker/cli/cli/command/formatter"
	h "github.com/lavrs/dms/pkg/client/http"
	"github.com/lavrs/dms/pkg/context"
	"github.com/lavrs/dms/pkg/logger"
)

// GetContainersLogs returns container logs
func GetContainersLogs(id string) (string, error) {
	logger.Debug("containers logs cmd")

	body, err := h.GET(context.Get().Address + "/api/logs/" + id)
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
func GetContainersMetrics(id string) ([]*formatter.ContainerStats, error) {
	logger.Debug("containers metrics cmd")

	body, err := h.GET(context.Get().Address + "/api/metrics/" + id)
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
func GetStoppedContainers() ([]string, error) {
	logger.Debug("stopped containers cmd")

	body, err := h.GET(context.Get().Address + "/api/stopped")
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
func GetLaunchedContainers() ([]string, error) {
	logger.Debug("launched containers cmd")

	body, err := h.GET(context.Get().Address + "/api/launched")
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
