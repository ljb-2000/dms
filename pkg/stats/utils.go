package stats

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/cli/command/formatter"
	"github.com/docker/docker/client"
	"math"
	"strings"
	"time"
)

func (s *Stats) collect(cli *client.Client, ID string) {
	s.changes.Lock()
	s.changes.changes[ID] = true
	s.changes.Unlock()

	for range time.Tick(time.Second) {
		stats, err := one(cli, ID)
		if err != nil {
			panic(err)
		}

		if stats.CPUPercentage == 0 {
			s.changes.Lock()
			s.changes.changes[ID] = false
			s.changes.Unlock()

			delete(s.data.data, ID)

			return
		}

		s.data.Lock()
		s.data.data[ID] = stats
		s.data.Unlock()
	}
}

func one(cli *client.Client, ID string) (*formatter.ContainerStats, error) {
	stats, err := cli.ContainerStats(context.Background(), ID, true)
	if err != nil {
		return nil, err
	}
	defer stats.Body.Close()

	dec := json.NewDecoder(stats.Body)
	var statsJSON *types.StatsJSON
	err = dec.Decode(&statsJSON)
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
