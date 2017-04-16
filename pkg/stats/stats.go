package stats

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/cli/command/formatter"
	"github.com/docker/docker/client"
	"strings"
	"sync"
	"time"
)

type dataMap struct {
	sync.RWMutex
	data map[string]*formatter.ContainerStats
}

type changeMap struct {
	sync.RWMutex
	changes map[string]bool
}

type Stats struct {
	data    dataMap
	changes changeMap
}

type statsAPI struct {
	Data     *[]*formatter.ContainerStats `json:"data,omitempty"`
	Launched []string                     `json:"launched,omitempty"`
	Stopped  []string                     `json:"stopped,omitempty"`
	Error    string                       `json:"error,omitempty"`
}

func (s *Stats) Collect() {
	s.data = dataMap{data: make(map[string]*formatter.ContainerStats)}
	s.changes = changeMap{changes: make(map[string]bool)}

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	for range time.Tick(time.Second * 3) {
		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
		if err != nil {
			panic(err)
		}

		for _, container := range containers {
			if _, ok := s.data.data[container.Names[0][1:]]; !ok {
				go s.collect(cli, container.Names[0][1:])
			}
		}
	}
}

func (s *Stats) Get(ID string) *statsAPI {
	var (
		stats      []*formatter.ContainerStats
		IDs        []string
		launched   []string
		stopped    []string
		isNotExist int = 0
	)

	if ID == "all" {
		s.data.RLock()
		for _, d := range s.data.data {
			IDs = append(IDs, d.Name)
		}
		s.data.RUnlock()
	} else if strings.Contains(ID, ",") {
		IDs = strings.Split(strings.Replace(ID, " ", "", -1), ",")
	} else {
		IDs = append(IDs, ID)
	}

	if len(s.changes.changes) != 0 {
		s.changes.RLock()
		for ID, status := range s.changes.changes {
			if status {
				launched = append(launched, ID)
			} else {
				stopped = append(stopped, ID)
			}
		}
		s.changes.RUnlock()

		for k := range s.changes.changes {
			delete(s.changes.changes, k)
		}
	}

	if len(s.data.data) == 0 {
		return &statsAPI{
			Launched: launched,
			Stopped:  stopped,
			Error:    "no running containers",
		}
	}

	s.data.RLock()
	for _, ID := range IDs {
		if data, ok := s.data.data[ID]; ok {
			stats = append(stats, data)
		} else {
			isNotExist++
		}
	}
	s.data.RUnlock()
	if isNotExist == len(IDs) {
		return &statsAPI{
			Launched: launched,
			Stopped:  stopped,
			Error:    "these containers are not running",
		}
	}

	return &statsAPI{
		Data:     &stats,
		Launched: launched,
		Stopped:  stopped,
	}
}
