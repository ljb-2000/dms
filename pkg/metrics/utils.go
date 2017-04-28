package metrics

import (
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/lavrs/docker-monitoring-service/pkg/docker"
	"github.com/lavrs/docker-monitoring-service/pkg/logger"
	"io"
	"time"
)

func (m *metrics) collect(id string) {
	m.changes.Lock()
	m.changes.changes[id] = true
	m.changes.Unlock()

	reader, err := docker.ContainerStats(id)
	if err != nil {
		logger.Panic(err)
	}
	defer reader.Close()

	dec := json.NewDecoder(reader)
	var statsJSON *types.StatsJSON

	for range time.Tick(m.uCMetricsInterval) {
		err = dec.Decode(&statsJSON)
		if err != nil {
			if err == io.EOF {
				logger.Info("container `", id, "` removed")
				m.removeCFromMap(id)
				return
			}
		}

		metrics := docker.Formatting(statsJSON)

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
