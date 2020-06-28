package canned

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/client"
	"github.com/gomodule/redigo/redis"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Tile38 struct {
	Container    testcontainers.Container
	DockerClient *client.Client

	Host   string
	Port   string
	Client redis.Conn
}

func NewTile38(ctx context.Context) (*Tile38, error) {
	os.Setenv("TC_HOST", "localhost")
	req := testcontainers.ContainerRequest{
		Image:        "tile38/tile38",
		ExposedPorts: []string{"9851/tcp"},
		WaitingFor:   wait.ForHTTP("/server").WithPort("9851"),
		AutoRemove:   true,
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, err
	}

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "9851")

	redisClient, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port.Port()))
	if err != nil {
		return nil, fmt.Errorf("error dialing redis, error: %v", err)
	}
	if err := redisClient.Send("OUTPUT", "json"); err != nil {
		return nil, fmt.Errorf("error setting output-json, error: %v", err)
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("error creating docker client, error: %v", err)
	}

	return &Tile38{
		Container:    container,
		DockerClient: dockerClient,

		Host:   host,
		Port:   port.Port(),
		Client: redisClient,
	}, nil
}
