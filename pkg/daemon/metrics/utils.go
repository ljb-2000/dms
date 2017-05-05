package metrics

import (
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/lavrs/dms/pkg/daemon/docker"
	"github.com/lavrs/dms/pkg/logger"
	"io"
	"time"
)

// collect metrics (container stats)
func (m *metrics) collect(id string) {
	// added to lunched containers array
	m.changes.Lock()
	m.changes.changes[id] = true
	m.changes.Unlock()

	// get container stats
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
			// container removed
			if err == io.EOF {
				logger.Info("container `", id, "` removed")
				m.removeCFromMaps(id)
				return
			}
			logger.Panic(err)
		}

		// formatting metrics
		metrics := docker.Formatting(statsJSON)

		// container stopped
		if metrics.CPUPercentage == 0 {
			logger.Info("container `", id, "` stopped")
			m.removeCFromMaps(id)
			return
		}

		// update metrics
		m.metrics.Lock()
		m.metrics.metrics[id] = metrics
		m.metrics.Unlock()
	}
}

// remove container from map
func (m *metrics) removeCFromMaps(id string) {
	// added to stopped containers array
	m.changes.Lock()
	m.changes.changes[id] = false
	m.changes.Unlock()

	// remove from metrics map
	m.metrics.Lock()
	delete(m.metrics.metrics, id)
	m.metrics.Unlock()
}
