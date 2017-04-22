package docker

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	c "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/cli/command/formatter"
	"github.com/docker/docker/client"
	"github.com/lavrs/docker-monitoring-service/pkg/logger"
	"io"
	"os"
    "time"
)

var (
	cli *client.Client
	err error
)

func init() {
	cli, err = client.NewEnvClient()
	if err != nil {
		logger.Panic(err)
	}
}

// Returns a list of running containers
func ContainerList() (*[]types.Container, error) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	return &containers, nil
}

// Returns the metric of the container
func ContainerStats(id string) (*formatter.ContainerStats, error) {
	stats, err := cli.ContainerStats(context.Background(), id, true)
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

	return Formatting(statsJSON), nil
}

// Create container
func ContainerCreate(cImage, cName string) error {
	_, err := cli.ContainerCreate(context.Background(), &c.Config{
		Image: cImage,
	}, &c.HostConfig{}, &network.NetworkingConfig{}, cName)
	if err != nil {
		return err
	}

	return nil
}

// Launches container
func ContainerStart(cName string) error {
	err = cli.ContainerStart(context.Background(), cName, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	return nil
}

// Download image
func ImagePull(cImage string) error {
	out, err := cli.ImagePull(context.Background(), cImage, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, out)
	if err != nil {
		return err
	}

	return nil
}

// Stops the container
func ContainerStop(cName string) error {
    t := time.Duration(0)
	err := cli.ContainerStop(context.Background(), cName, &t)
	if err != nil {
		return err
	}

	return nil
}

// Removes image
func ImageRemove(cImage string) error {
	_, err := cli.ImageRemove(context.Background(), cImage, types.ImageRemoveOptions{})
	if err != nil {
		return err
	}

	return nil
}

// Removes container
func ContainerRemove(cName string) error {
	err := cli.ContainerRemove(context.Background(), cName, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}

	return nil
}
