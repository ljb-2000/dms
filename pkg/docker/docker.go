package docker

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	c "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
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

func StartContainer(t *testing.T, cImage, cName string) {
	out, err := cli.ImagePull(ctx, cImage, types.ImagePullOptions{})
	assert.NoError(t, err)
	io.Copy(os.Stdout, out)
	assert.NoError(t, err)

	container, err := cli.ContainerCreate(ctx, &c.Config{
		Image: cImage,
	}, &c.HostConfig{}, &network.NetworkingConfig{}, cName)
	assert.NoError(t, err)

	err = cli.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})
	assert.NoError(t, err)
}

func StopContainer(t *testing.T, cName string) {
	err := cli.ContainerStop(ctx, cName, nil)
	assert.NoError(t, err)
}

func RemoveContainer(t *testing.T, cName string) {
	err := cli.ContainerRemove(ctx, cName, types.ContainerRemoveOptions{})
	assert.NoError(t, err)

	// проверять имадж и удалять
}
