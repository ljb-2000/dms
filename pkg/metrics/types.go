package metrics

import (
	"github.com/docker/docker/cli/command/formatter"
	"sync"
)

type dataMap struct {
	sync.RWMutex
	data map[string]*formatter.ContainerStats
}

type changeMap struct {
	sync.RWMutex
	changes map[string]bool
}

type metrics struct {
	data    dataMap
	changes changeMap
}

type metricsAPI struct {
	Metrics  []*formatter.ContainerStats `json:"metrics,omitempty"`
	Launched []string                    `json:"launched,omitempty"`
	Stopped  []string                    `json:"stopped,omitempty"`
	Message  string                      `json:"message,omitempty"`
}
