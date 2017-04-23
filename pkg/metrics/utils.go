package metrics

import (
	"github.com/docker/docker/cli/command/formatter"
	"github.com/lavrs/docker-monitoring-service/pkg/docker"
	"github.com/lavrs/docker-monitoring-service/pkg/logger"
	"io"
	"time"
)

func (m *metrics) collect(id string) {
	m.changes.Lock()
	m.changes.changes[id] = true
	m.changes.Unlock()

	for range time.Tick(m.uCMetricsInterval) {
		metrics, err := one(id)
		if err != nil {
			if err == io.EOF {
				logger.Info("container `", id, "` removed")
				m.removeCFromMap(id)
				return
			}
			logger.Panic(err)
		}

		if metrics.CPUPercentage == 0 {
			logger.Info("container `", id, "` stopped")
			m.removeCFromMap(id)
			return
		}

		m.metrics.Lock()
		m.metrics.metrics[id] = metrics
		m.metrics.Unlock()
	}
}

func (m *metrics) removeCFromMap(id string) {
	m.changes.Lock()
	m.changes.changes[id] = false
	m.changes.Unlock()

	m.metrics.Lock()
	delete(m.metrics.metrics, id)
	m.metrics.Unlock()
}

func one(id string) (*formatter.ContainerStats, error) {
	statsJSON, err := docker.ContainerStats(id)
	if err != nil {
		return nil, err
	}

	return statsJSON, nil
}
