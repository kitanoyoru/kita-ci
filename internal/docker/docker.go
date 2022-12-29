package docker

import (
	"bytes"
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"github.com/kitanoyoru/kita-ci/pkg/config"
	"github.com/kitanoyoru/kita-ci/pkg/structs"
)

type DockerClient struct {
	cli *client.Client
	cfg config.DockerConfig
}

func NewDockerClient(config config.DockerConfig) *DockerClient {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Panicf("Failed to establish connection with host Docker server: %d", err)
	}

	return &DockerClient{
		cli: cli,
		cfg: config,
	}
}

func (c *DockerClient) RunCIBuilderContainer(payload structs.CIBuilderPayload) (string, error) {
	resp, err := c.cli.ContainerCreate(context.Background(), &container.Config{
		Image: c.cfg.Image,
		Env:   []string{payload.RepoURL, payload.Branch, payload.Tag},
		Volumes: map[string]struct{}{
			"/var/run/docker.sock": "/var/run/docker.sock",
		},
	}, nil, nil, nil, "")
	if err != nil {
		return "", err
	}

	if err := c.cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	statusCh, errCh := c.cli.ContainerWait(context.Background(), resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", err
		}
	case <-statusCh:
	}

	reader, err := c.cli.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return "", err
	}

	// REFACTOR: Move to utils
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(reader)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
