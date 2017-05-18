package metrics

import (
	"github.com/docker/cli/cli/command/formatter"
	"github.com/lavrs/dms/pkg/daemon/docker"
	"github.com/lavrs/dms/pkg/logger"
	"strings"
	"time"
)

// init metrics
var m = &metrics{
	metrics:             metricsMap{metrics: make(map[string]*formatter.ContainerStats)},
	changes:             changeMap{changes: make(map[string]bool)},
	updCListInterval:    time.Second * 3,
	updCMetricsInterval: time.Second * 1,
}

// Get returns metrics obj
func Get() *metrics {
	return m
}

// SetUCMetricsInterval set update container metrics interval
func (m *metrics) SetUCMetricsInterval(t time.Duration) {
	logger.Info("set update container metrics interval =", t)
	m.updCMetricsInterval = t
}

// SetUCListInterval set update containers list interval
func (m *metrics) SetUCListInterval(t time.Duration) {
	logger.Info("set update containers list interval =", t)
	m.updCListInterval = t
}

// Collect collect metrics (check new containers)
func (m *metrics) Collect() {
	// if metrics already collects, returns
	if m.started {
		return
	}
	m.started = true

	for range time.Tick(m.updCListInterval) {
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

// GetContainerLogs returns container logs
func GetContainerLogs(id string) string {
	logs, err := docker.ContainersLogs(id)
	if err != nil {
		return err.Error()
	}

	logger.Info(id, "container logs:", logs)
	return logs
}

// GetStoppedContainers returns stopped containers
func (m *metrics) GetStoppedContainers() []string {
	var (
		stopped   []string
		isStopped = false
	)

	// parse changes
	if len(m.changes.changes) != 0 {
		m.changes.RLock()
		for id, status := range m.changes.changes {
			if !status {
				isStopped = true
				stopped = append(stopped, id)
			}
		}
		m.changes.RUnlock()
	}
	logger.Info("stopped containers:", stopped)

	// there are stopped containers
	if isStopped {
		return stopped
	}

	// no stopped containers
	return []string{"no stopped containers"}
}

// GetLaunchedContainers returns Launched containers
func (m *metrics) GetLaunchedContainers() []string {
	var (
		launched   []string
		isLaunched = false
	)

	// parse changes
	m.changes.RLock()
	if len(m.changes.changes) != 0 {
		for id, status := range m.changes.changes {
			if status {
				isLaunched = true
				launched = append(launched, id)
			}
		}
	} else {
		logger.Info("no changes")
	}
	m.changes.RUnlock()
	logger.Info("launched containers:", launched)

	// there are launched containers
	if isLaunched {
		return launched
	}

	// no launched containers
	return []string{"no launched containers"}
}

// Get returns container(s) metrics
func (m *metrics) Get(id string) *metricsAPI {
	var (
		metrics                []*formatter.ContainerStats
		ids, launched, stopped []string
		isNotExist             = 0
	)

	// parse id (all / one / ... containers)
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
	logger.Info("get", ids, "container(s) metrics")

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

		// flush changes
		for k := range m.changes.changes {
			delete(m.changes.changes, k)
		}
		m.changes.RUnlock()
	}
	logger.Info("stopped:", stopped, "\nlaunched:", launched)

	// return if no running containers
	if len(m.metrics.metrics) == 0 {
        logger.Info("no running containers")
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
			continue
		}
		// if container are not running
		logger.Info("container `", id, "` are not running")
		isNotExist++
	}
	m.metrics.RUnlock()
	// returns if all specified containers are not running
	if isNotExist == len(ids) {
		logger.Info("these containers:", ids, "are not running")
		return &metricsAPI{
			Launched: launched,
			Stopped:  stopped,
			Message:  "these containers are not running",
		}
	}

	// returns metrics
	logger.Info(ids, "metrics: ", metrics)
	return &metricsAPI{
		Metrics:  &metrics,
		Launched: launched,
		Stopped:  stopped,
	}
}
