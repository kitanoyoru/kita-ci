package temp

import (
	"context"
	"io"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"github.com/kitanoyoru/kita-ci/pkg/config"
)

type DockerClient struct {
	cli *client.Client
	cfg *config.DockerConfig
}

func NewDockerClient(config *config.DockerConfig) *DockerClient {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Panicf("Failed to establish connection with host Docker server: %d", err)
	}

	return &DockerClient{
		cli: cli,
		cfg: config,
	}
}

func (c *DockerClient) RunCIBuilderContainer() io.ReadCloser {
	resp, err := c.cli.ContainerCreate(context.Background(), &container.Config{
		Image: c.cfg.Image,
	}, nil, nil, nil, "")
	if err != nil {
		log.Panicf("Failed to create container: %d", err)
	}

	if err := c.cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Panicf("Failed to start container: %d", err)
	}

	statusCh, errCh := c.cli.ContainerWait(context.Background(), resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			log.Panicf("Error in running container: %d", err)
		}
	case <-statusCh:
	}

	out, err := c.cli.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		log.Panicf("Failed to get logs from container: %d", err)
	}

	return out
}
