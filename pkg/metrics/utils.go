package metrics

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/cli/command/formatter"
	"github.com/lavrs/docker-monitoring-service/pkg/docker"
	"math"
	"strings"
	"time"
)

func (m *metrics) collect(id string) {
	m.changes.Lock()
	m.changes.changes[id] = true
	m.changes.Unlock()

	for range time.Tick(m.ucTime) {
		metrics, err := one(id)
		if err != nil {
			panic(err)
		}

		if metrics.CPUPercentage == 0 {
			m.changes.Lock()
			m.changes.changes[id] = false
			m.changes.Unlock()

			delete(m.data.data, id)

			return
		}

		m.data.Lock()
		m.data.data[id] = metrics
		m.data.Unlock()
	}
}

func one(id string) (*formatter.ContainerStats, error) {
	statsJSON, err := docker.ContainerStats(id)
	if err != nil {
		return nil, err
	}

	return docker.Formatting(statsJSON), nil
}
