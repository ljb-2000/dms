package metrics

import (
	"github.com/docker/docker/cli/command/formatter"
	"github.com/lavrs/docker-monitoring-service/pkg/docker"
	"github.com/lavrs/docker-monitoring-service/pkg/logger"
	"github.com/pkg/errors"
	"strings"
	"time"
)

var m metrics

// Get metrics
func Get() *metrics {
	return &m
}

// Create new metrics
func NewMetrics() (*metrics, error) {
	logger.Info("new metrics")

	if m.isCreated {
		return nil, errors.New("metrics already create")
	}

	m.isCreated = true
	m.data = metricsMap{metrics: make(map[string]*formatter.ContainerStats)}
	m.changes = changeMap{changes: make(map[string]bool)}
	m.ucListTime = time.Second * 3
	m.ucTime = time.Second

	return &m, nil
}

func (m *metrics) SetUCLTime(t time.Duration) {
	m.ucListTime = t
}

func (m *metrics) SetUCTime(t time.Duration) {
	m.ucTime = t
}

func (m *metrics) Collect() {
    logger.Info(m.ucListTime)
	for range time.Tick(m.ucListTime) {
		containers, err := docker.ContainerList()
		if err != nil {
			logger.Panic(err)
		}
        logger.Info(containers)

		for _, container := range *containers {
			if _, ok := m.data.metrics[container.Names[0][1:]]; !ok {
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
		m.data.RLock()
		for _, d := range m.data.metrics {
			ids = append(ids, d.Name)
		}
		m.data.RUnlock()
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
		m.changes.RUnlock()

		for k := range m.changes.changes {
			delete(m.changes.changes, k)
		}
	} else {
		logger.Info("no new containers")
	}

	if len(m.data.metrics) == 0 {
		logger.Info("no running container")
		return &metricsAPI{
			Launched: launched,
			Stopped:  stopped,
			Message:  "no running containers",
		}
	}

	m.data.RLock()
	for _, id := range ids {
		if data, ok := m.data.metrics[id]; ok {
			metrics = append(metrics, data)
		} else {
			logger.Info("container `", id, "` are not running")
			isNotExist++
		}
	}
	m.data.RUnlock()
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
