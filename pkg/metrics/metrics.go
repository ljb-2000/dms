package metrics

import (
	"github.com/docker/docker/cli/command/formatter"
	"github.com/lavrs/docker-monitoring-service/pkg/docker"
	"github.com/lavrs/docker-monitoring-service/pkg/logger"
	"strings"
	"time"
)

var m = &metrics{
	metrics:           metricsMap{metrics: make(map[string]*formatter.ContainerStats)},
	changes:           changeMap{changes: make(map[string]bool)},
	uCListInterval:    time.Second * 3,
	uCMetricsInterval: time.Second * 1,
}

// Get metrics
func Get() *metrics {
	return m
}

func (m *metrics) SetUCMetricsInterval(t time.Duration) {
	m.uCMetricsInterval = t
}

func (m *metrics) SetUCListInterval(t time.Duration) {
	m.uCListInterval = t
}

func (m *metrics) Collect() {
	for range time.Tick(m.uCListInterval) {
		containers, err := docker.ContainerList()
		if err != nil {
			logger.Panic(err)
		}

		for _, container := range *containers {
			if _, ok := m.metrics.metrics[container.Names[0][1:]]; !ok {
				logger.Info("new container `", container.Names[0][1:], "`")

				go m.collect(container.Names[0][1:])
			}
		}
	}
}

func (m *metrics) Get(id string) *metricsAPI {
	logger.Info("get container(s) metrics")

	var (
		metrics    []*formatter.ContainerStats
		ids        []string
		launched   []string
		stopped    []string
		isNotExist = 0
	)

	if id == "all" {
		m.metrics.RLock()
		for _, d := range m.metrics.metrics {
			ids = append(ids, d.Name)
		}
		m.metrics.RUnlock()
	} else if strings.Contains(id, " ") {
		ids = strings.Split(strings.Replace(id, " ", "", -1), " ")
	} else {
		ids = append(ids, id)
	}

	if len(m.changes.changes) != 0 {
		m.changes.RLock()
		for id, status := range m.changes.changes {
			if status {
				launched = append(launched, id)
			} else {
				stopped = append(stopped, id)
			}
		}

		for k := range m.changes.changes {
			delete(m.changes.changes, k)
		}
		m.changes.RUnlock()
	} else {
		logger.Info("no new containers")
	}

	if len(m.metrics.metrics) == 0 {
		logger.Info("no running container")
		return &metricsAPI{
			Launched: launched,
			Stopped:  stopped,
			Message:  "no running containers",
		}
	}

	m.metrics.RLock()
	for _, id := range ids {
		if data, ok := m.metrics.metrics[id]; ok {
			metrics = append(metrics, data)
		} else {
			logger.Info("container `", id, "` are not running")
			isNotExist++
		}
	}
	m.metrics.RUnlock()
	if isNotExist == len(ids) {
		logger.Info("these containers are not running")
		return &metricsAPI{
			Launched: launched,
			Stopped:  stopped,
			Message:  "these containers are not running",
		}
	}

	logger.Info("return metrics")
	return &metricsAPI{
		Metrics:  &metrics,
		Launched: launched,
		Stopped:  stopped,
	}
}
