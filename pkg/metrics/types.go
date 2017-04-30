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

// type for collect metrics
type metrics struct {
	metrics           metricsMap
	changes           changeMap
	uCListInterval    time.Duration
	uCMetricsInterval time.Duration
}

// type for API
type metricsAPI struct {
	Metrics  *[]*formatter.ContainerStats `json:"metrics,omitempty"`
	Launched []string                     `json:"launched,omitempty"`
	Stopped  []string                     `json:"stopped,omitempty"`
	Message  string                       `json:"message,omitempty"`
}
