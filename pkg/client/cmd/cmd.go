package cmd

import (
	"encoding/json"
	"errors"
	"github.com/docker/docker/cli/command/formatter"
	"github.com/lavrs/dms/pkg/logger"
	"io/ioutil"
	"net/http"
)

// ContainersLogs show container logs
func ContainersLogsCmd() {
}

// ContainersMetrics show container(s) metrics
func ContainersMetricsCmd() {
}

// StoppedContainers stopped containers
func StoppedContainersCmd() {
}

// LaunchedContainers launched containers
func LaunchedContainersCmd() {
}

// GetContainersLogs returns container logs
func GetContainersLogs(url, id string) (string, error) {
	logger.Info("containers logs cmd")

	resp, err := http.Get(url + "/logs/" + id)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
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

	resp, err := http.Get(url + "/metrics/" + id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
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

	resp, err := http.Get(url + "/stopped")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
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

	resp, err := http.Get(url + "/launched")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
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
