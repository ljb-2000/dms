package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	c "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/lavrs/dms/pkg/logger"
	"io"
	"io/ioutil"
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

// ContainerList returns a list of running containers
func ContainerList() (*[]types.Container, error) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	return &containers, nil
}

// ContainerStats returns container metrics channel
func ContainerStats(id string) (io.ReadCloser, error) {
	cStats, err := cli.ContainerStats(context.Background(), id, true)
	if err != nil {
		return nil, err
	}

	return cStats.Body, nil
}

// ContainerCreate create container
func ContainerCreate(cImage, cName string) error {
	_, err := cli.ContainerCreate(context.Background(), &c.Config{
		Image: cImage,
	}, &c.HostConfig{}, &network.NetworkingConfig{}, cName)
	if err != nil {
		return err
	}

	return nil
}

// ContainerStart launches container
func ContainerStart(cName string) error {
	err = cli.ContainerStart(context.Background(), cName, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	return nil
}

// ImagePull download image
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

// ContainerStop stops the container
func ContainerStop(cName string) error {
	t := time.Duration(0)
	err := cli.ContainerStop(context.Background(), cName, &t)
	if err != nil {
		return err
	}

	return nil
}

// ImageRemove removes image
func ImageRemove(cImage string) error {
	_, err := cli.ImageRemove(context.Background(), cImage, types.ImageRemoveOptions{})
	if err != nil {
		return err
	}

	return nil
}

// ContainerRemove removes container
func ContainerRemove(cName string) error {
	err := cli.ContainerRemove(context.Background(), cName, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
	if err != nil {
		return err
	}

	return nil
}

// ContainersLogs returns containers logs
func ContainersLogs(cName string) (string, error) {
	reader, err := cli.ContainerLogs(context.Background(), cName, types.ContainerLogsOptions{
		ShowStdout: true,
		Since:      time.Since(time.Now().AddDate(0, 0, -1)).String(),
	})
	if err != nil {
		return "", err
	}
	defer reader.Close()

	logs, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(logs), nil
}
