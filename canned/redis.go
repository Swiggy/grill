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

type Redis struct {
	Container    testcontainers.Container
	DockerClient *client.Client

	Host string
	Port string
	Pool *redis.Pool
}

func NewRedis(ctx context.Context) (*Redis, error) {
	os.Setenv("TC_HOST", "localhost")

	req := testcontainers.ContainerRequest{
		Image:        getEnvString("REDIS_CONTAINER_IMAGE", "redis:7.0.11"),
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379"),
		AutoRemove:   true,
		SkipReaper:   skipReaper(),
		RegistryCred: getBasicAuth(),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, err
	}

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "6379")

	redisPool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port.Port()))
		},
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("error creating docker client, error: %v", err)
	}

	return &Redis{
		Container:    container,
		DockerClient: dockerClient,

		Host: host,
		Port: port.Port(),
		Pool: redisPool,
	}, nil
}
