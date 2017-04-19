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

	for range time.Tick(time.Second) {
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

	return formatting(statsJSON), nil
}

func formatting(statsJSON *types.StatsJSON) *formatter.ContainerStats {
	var stats formatter.ContainerStats

	read, write := parseBlockIO(statsJSON.BlkioStats)
	received, sent := parseNetwork(statsJSON.Networks)

	stats.SetStatistics(formatter.StatsEntry{
		Name:             statsJSON.Name[1:],
		ID:               statsJSON.ID,
		MemoryPercentage: parseMemory(statsJSON),
		CPUPercentage:    parseCPU(statsJSON),
		Memory:           float64(statsJSON.MemoryStats.Usage),
		MemoryLimit:      float64(statsJSON.MemoryStats.Limit),
		NetworkTx:        sent,
		NetworkRx:        received,
		BlockRead:        float64(read),
		BlockWrite:       float64(write),
		PidsCurrent:      statsJSON.PidsStats.Current,
	})

	return &stats
}

func parseMemory(statsJSON *types.StatsJSON) float64 {
	mem := float64(statsJSON.MemoryStats.Usage) / float64(statsJSON.MemoryStats.Limit) * 100.0

	if math.IsNaN(mem) {
		return 0
	}
	return mem
}

func parseNetwork(net map[string]types.NetworkStats) (float64, float64) {
	var rx, tx float64

	for _, v := range net {
		rx += float64(v.RxBytes)
		tx += float64(v.TxBytes)
	}

	return rx, tx
}

func parseCPU(statsJSON *types.StatsJSON) float64 {
	var (
		cpuPercent = 0.0

		cpuDelta = float64(statsJSON.CPUStats.CPUUsage.TotalUsage) - float64(statsJSON.PreCPUStats.CPUUsage.TotalUsage)

		systemDelta = float64(statsJSON.CPUStats.SystemUsage) - float64(statsJSON.PreCPUStats.SystemUsage)
		onlineCPUs  = float64(statsJSON.CPUStats.OnlineCPUs)
	)

	if onlineCPUs == 0.0 {
		onlineCPUs = float64(len(statsJSON.CPUStats.CPUUsage.PercpuUsage))
	}
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * onlineCPUs * 100.0
	}
	return cpuPercent
}

func parseBlockIO(blkioStats types.BlkioStats) (uint64, uint64) {
	var (
		blkRead  uint64
		blkWrite uint64
	)

	for _, blkIO := range blkioStats.IoServiceBytesRecursive {
		switch strings.ToLower(blkIO.Op) {
		case "read":
			blkRead = blkRead + blkIO.Value
		case "write":
			blkWrite = blkWrite + blkIO.Value
		}
	}

	return blkRead, blkWrite
}
