package metrics

import (
	"github.com/docker/docker/cli/command/formatter"
	"github.com/lavrs/docker-monitoring-service/pkg/docker"
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
				return
			}

			panic(err)
		}

		if metrics.CPUPercentage == 0 {
			m.changes.Lock()
			m.changes.changes[id] = false
			m.changes.Unlock()

			m.data.Lock()
			delete(m.data.data, id)
			m.data.Unlock()

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
