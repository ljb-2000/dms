package metrics

import (
	"github.com/docker/docker/cli/command/formatter"
	"github.com/lavrs/docker-monitoring-service/pkg/docker"
	"strings"
	"time"
)

func NewMetrics() *metrics {
	return &metrics{
		data:       dataMap{data: make(map[string]*formatter.ContainerStats)},
		changes:    changeMap{changes: make(map[string]bool)},
		ucListTime: time.Second * 3,
		ucTime:     time.Second,
	}
}

func (m *metrics) SetUCLTime(t time.Duration) {
	m.ucListTime = t
}

func (m *metrics) SetUCTime(t time.Duration) {
	m.ucTime = t
}

func (m *metrics) Collect() {
	for range time.Tick(m.ucListTime) {
		containers, err := docker.ContainerList()
		if err != nil {
			panic(err)
		}

		for _, container := range *containers {
			if _, ok := m.data.data[container.Names[0][1:]]; !ok {
				go m.collect(container.Names[0][1:])
			}
		}
	}
}

func (m *metrics) Get(id string) *metricsAPI {
	var (
		metrics    []*formatter.ContainerStats
		ids        []string
		launched   []string
		stopped    []string
		isNotExist = 0
	)

	if id == "all" {
		m.data.RLock()
		for _, d := range m.data.data {
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
	}

	if len(m.data.data) == 0 {
		return &metricsAPI{
			Launched: launched,
			Stopped:  stopped,
			Message:  "no running containers",
		}
	}

	m.data.RLock()
	for _, id := range ids {
		if data, ok := m.data.data[id]; ok {
			metrics = append(metrics, data)
		} else {
			isNotExist++
		}
	}
	m.data.RUnlock()
	if isNotExist == len(ids) {
		return &metricsAPI{
			Launched: launched,
			Stopped:  stopped,
			Message:  "these containers are not running",
		}
	}

	return &metricsAPI{
		Metrics:  &metrics,
		Launched: launched,
		Stopped:  stopped,
	}
}
