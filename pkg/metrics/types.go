package metrics

import (
	"github.com/docker/docker/cli/command/formatter"
	"sync"
	"time"
)

type metricsMap struct {
	sync.RWMutex
	metrics map[string]*formatter.ContainerStats
}

type changeMap struct {
	sync.RWMutex
	changes map[string]bool
}

type metrics struct {
	data       metricsMap
	changes    changeMap
	ucListTime time.Duration
	ucTime     time.Duration
	isCreated  bool
}

type metricsAPI struct {
	Metrics  *[]*formatter.ContainerStats `json:"metrics,omitempty"`
	Launched []string                     `json:"launched,omitempty"`
	Stopped  []string                     `json:"stopped,omitempty"`
	Message  string                       `json:"message,omitempty"`
}
