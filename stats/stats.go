package stats

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/cli/command/formatter"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"math"
	"strings"
)

func getMoreThanOne(cli *client.Client, IDs []string) (*[]*formatter.ContainerStats, error) {
	var res []*formatter.ContainerStats

	for _, ID := range IDs {
		oneStats, err := getOne(cli, ID)
		if err != nil {
			return nil, err
		}

		res = append(res, oneStats)
	}

	return &res, nil
}

func GetStats(cli *client.Client, ID string) (*[]*formatter.ContainerStats, error) {
	if ID == "all" {
		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
		if err != nil {
			return nil, err
		}
		if len(containers) == 0 {
			return nil, errors.New("no running containers")
		}

		var IDs []string
		for _, container := range containers {
			IDs = append(IDs, container.ID)
		}
		res, err := getMoreThanOne(cli, IDs)
		if err != nil {
			return nil, err
		}

		return res, nil
	} else if strings.Contains(ID, ",") {
		IDs := strings.Split(strings.Replace(ID, " ", "", -1), ",")
		res, err := getMoreThanOne(cli, IDs)
		if err != nil {
			return nil, err
		}

		return res, nil
	} else {
		var res []*formatter.ContainerStats

		oneStat, err := getOne(cli, ID)
		if err != nil {
			return nil, err
		}
		res = append(res, oneStat)

		return &res, nil
	}
}

func getOne(cli *client.Client, ID string) (*formatter.ContainerStats, error) {
	res, err := cli.ContainerStats(context.Background(), ID, true)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	dec := json.NewDecoder(res.Body)
	var resJSON *types.StatsJSON
	err = dec.Decode(&resJSON)
	if err != nil {
		return nil, err
	}

	return convert(resJSON), nil
}

func convert(statsJSON *types.StatsJSON) *formatter.ContainerStats {
	var containerStats formatter.ContainerStats

	read, write := parseBlockIO(statsJSON.BlkioStats)
	received, sent := parseNetwork(statsJSON.Networks)

	containerStats.SetStatistics(formatter.StatsEntry{
		Name:             statsJSON.Name[1:],
		ID:               statsJSON.ID,
		MemoryPercentage: memoryToPercentage(statsJSON),
		CPUPercentage:    parseCPU(statsJSON),
		Memory:           float64(statsJSON.MemoryStats.Usage),
		MemoryLimit:      float64(statsJSON.MemoryStats.Limit),
		NetworkTx:        sent,
		NetworkRx:        received,
		BlockRead:        float64(read),
		BlockWrite:       float64(write),
		PidsCurrent:      statsJSON.PidsStats.Current,
	})

	return &containerStats
}

func memoryToPercentage(statsJSON *types.StatsJSON) float64 {
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
