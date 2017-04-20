package docker

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	c "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"io"
	"os"
)

var (
	ctx = context.Background()
	cli *client.Client
	err error
)

func init() {
	cli, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}
}

func ContainerList() ([]types.Container, error) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	return containers, nil
}

func ContainerStats(id string) (*types.StatsJSON, error) {
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

	return statsJSON, nil
}

func ContainerCreate(cImage, cName string) error {
	_, err := cli.ContainerCreate(ctx, &c.Config{
		Image: cImage,
	}, &c.HostConfig{}, &network.NetworkingConfig{}, cName)
	if err != nil {
		return err
	}

	return nil
}

func ContainerStart(cName string) error {
	err = cli.ContainerStart(ctx, cName, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	return nil
}

func ImagePull(cImage string) error {
	out, err := cli.ImagePull(ctx, cImage, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, out)
	if err != nil {
		return err
	}

	return nil
}

func ContainerStop(cName string) error {
	err := cli.ContainerStop(ctx, cName, nil)
	if err != nil {
		return err
	}

	return nil
}

func ImageRemove(cImage string) error {
	_, err := cli.ImageRemove(ctx, cImage, types.ImageRemoveOptions{})
	if err != nil {
		return err
	}

	return nil
}

func ContainerRemove(cName string) error {
	err := cli.ContainerRemove(ctx, cName, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}

	return nil
}
