package canned

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/client"
	"github.com/gomodule/redigo/redis"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Tile38 struct {
	Container    testcontainers.Container
	DockerClient *client.Client

	Host string
	Port string
	Pool *redis.Pool
}

func NewTile38(ctx context.Context) (*Tile38, error) {
	os.Setenv("TC_HOST", "localhost")

	req := testcontainers.ContainerRequest{
		Image:        getEnvString("TILE38_CONTAINER_IMAGE", "tile38/tile38:1.31.0"),
		SkipReaper:   skipReaper(),
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

	redisPool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port.Port()))
			if err != nil {
				return nil, err
			}
			if err := conn.Send("OUTPUT", "json"); err != nil {
				return nil, fmt.Errorf("error setting output-json, error: %v", err)
			}
			return conn, nil
		},
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("error creating docker client, error: %v", err)
	}

	return &Tile38{
		Container:    container,
		DockerClient: dockerClient,

		Host: host,
		Port: port.Port(),
		Pool: redisPool,
	}, nil
}
