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
	body, err := h.GET(context.Get().Address + "/api/logs/" + id)
	if err != nil {
		return "", err
	}

	var logs struct{ Logs string }
	err = json.Unmarshal(body, &logs)
	if err != nil {
		return "", err
	}

	logger.Info(id, "container logs:", logs.Logs)
	return logs.Logs, nil
}

// GetContainersMetrics returns container(s) metrics
func GetContainersMetrics(id string) ([]*formatter.ContainerStats, error) {
	body, err := h.GET(context.Get().Address + "/api/metrics/" + id)
	if err != nil {
		return nil, err
	}

	var metrics struct{ Metrics []*formatter.ContainerStats }
	err = json.Unmarshal(body, &metrics)
	if err != nil {
		return nil, err
	}

	logger.Info(id, "container(s) metrics:", metrics.Metrics)
	return metrics.Metrics, nil
}

// GetStoppedContainers returns stopped containers
func GetStoppedContainers() ([]string, error) {
	body, err := h.GET(context.Get().Address + "/api/stopped")
	if err != nil {
		return nil, err
	}

	var stopped struct{ Stopped []string }
	err = json.Unmarshal(body, &stopped)
	if err != nil {
		return nil, err
	}

	logger.Info("stopped containers:", stopped.Stopped)
	return stopped.Stopped, nil
}

// GetLaunchedContainers returns launched containers
func GetLaunchedContainers() ([]string, error) {
	body, err := h.GET(context.Get().Address + "/api/launched")
	if err != nil {
		return nil, err
	}

	var launched struct{ Launched []string }
	err = json.Unmarshal(body, &launched)
	if err != nil {
		return nil, err
	}

	logger.Info("launched containers:", launched.Launched)
	return launched.Launched, nil
}
