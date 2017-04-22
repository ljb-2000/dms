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

	for range time.Tick(m.ucTime) {
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

		m.data.Lock()
		m.data.metrics[id] = metrics
		m.data.Unlock()
	}
}

func (m *metrics) removeCFromMap(id string) {
	m.changes.Lock()
	m.changes.changes[id] = false
	m.changes.Unlock()

	m.data.Lock()
	delete(m.data.metrics, id)
	m.data.Unlock()
}

func one(id string) (*formatter.ContainerStats, error) {
	statsJSON, err := docker.ContainerStats(id)
	if err != nil {
		return nil, err
	}

	return statsJSON, nil
}
