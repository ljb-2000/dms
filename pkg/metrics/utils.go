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

	var (
		metrics *formatter.ContainerStats
		data    = make(chan *formatter.ContainerStats)
		done    = make(chan bool)
		err     error
	)

	go func() {
		err = docker.ContainerStats(id, data, done)
		if err != nil {
			if err == io.EOF {
				logger.Info("container `", id, "` removed")
				m.removeCFromMap(id, data, done)
				return
			}
			logger.Panic(err)
		}
	}()

	for range time.Tick(m.uCMetricsInterval) {
		metrics = <-data

		if metrics.CPUPercentage == 0 {
			logger.Info("container `", id, "` stopped")
			m.removeCFromMap(id, data, done)
			return
		}

		m.metrics.Lock()
		m.metrics.metrics[id] = metrics
		m.metrics.Unlock()
	}
}

func (m *metrics) removeCFromMap(id string, data chan *formatter.ContainerStats, done chan bool) {
	m.changes.Lock()
	m.changes.changes[id] = false
	m.changes.Unlock()

	m.metrics.Lock()
	delete(m.metrics.metrics, id)
	m.metrics.Unlock()

	close(done)
	close(data)
}
