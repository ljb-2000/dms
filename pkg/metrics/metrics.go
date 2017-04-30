package metrics

import (
	"github.com/docker/docker/cli/command/formatter"
	"github.com/lavrs/docker-monitoring-service/pkg/docker"
	"github.com/lavrs/docker-monitoring-service/pkg/logger"
	"strings"
	"time"
)

// init metrics
var m = &metrics{
	metrics:           metricsMap{metrics: make(map[string]*formatter.ContainerStats)},
	changes:           changeMap{changes: make(map[string]bool)},
	uCListInterval:    time.Second * 3,
	uCMetricsInterval: time.Second * 1,
}

// get metrics obj
func Get() *metrics {
	return m
}

// set update container metrics interval
func (m *metrics) SetUCMetricsInterval(t time.Duration) {
	m.uCMetricsInterval = t
}

// set update containers list interval
func (m *metrics) SetUCListInterval(t time.Duration) {
	m.uCListInterval = t
}

// collect metrics (check new containers)
func (m *metrics) Collect() {
	for range time.Tick(m.uCListInterval) {
		containers, err := docker.ContainerList()
		if err != nil {
			logger.Panic(err)
		}

		for _, container := range *containers {
			if _, ok := m.metrics.metrics[container.Names[0][1:]]; !ok {
				logger.Info("new container `", container.Names[0][1:], "`")

                // start collect new metrics
				go m.collect(container.Names[0][1:])
			}
		}
	}
}

// returns container(s) metrics
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

	// parse changes
	if len(m.changes.changes) != 0 {
		m.changes.RLock()
		for id, status := range m.changes.changes {
			if status {
				launched = append(launched, id)
			} else {
				stopped = append(stopped, id)
			}
		}

		// clear changes
		for k := range m.changes.changes {
			delete(m.changes.changes, k)
		}
		m.changes.RUnlock()
	} else {
		logger.Info("no new containers")
	}

	// return if no running containers
	if len(m.metrics.metrics) == 0 {
		logger.Info("no running container")
		return &metricsAPI{
			Launched: launched,
			Stopped:  stopped,
			Message:  "no running containers",
		}
	}

	// get containers metrics from data map
	m.metrics.RLock()
	for _, id := range ids {
		if data, ok := m.metrics.metrics[id]; ok {
			metrics = append(metrics, data)
		} else {
			// if container are not running
			logger.Info("container `", id, "` are not running")
			isNotExist++
		}
	}
	m.metrics.RUnlock()
	// returns if all specified containers are not running
	if isNotExist == len(ids) {
		logger.Info("these containers are not running")
		return &metricsAPI{
			Launched: launched,
			Stopped:  stopped,
			Message:  "these containers are not running",
		}
	}

	// returns metrics
	logger.Info("return metrics")
	return &metricsAPI{
		Metrics:  &metrics,
		Launched: launched,
		Stopped:  stopped,
	}
}
