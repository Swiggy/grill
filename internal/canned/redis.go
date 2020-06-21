package canned

import (
	"context"
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Redis struct {
	Container testcontainers.Container

	Host   string
	Port   string
	Client redis.Conn
}

func NewRedis(ctx context.Context) (*Redis, error) {
	req := testcontainers.ContainerRequest{
		Image:        "redis",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379"),
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

	client, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port.Port()))
	if err != nil {
		return nil, fmt.Errorf("error dialing redis, error: %v", err)
	}

	return &Redis{
		Container: container,
		Host:      host,
		Port:      port.Port(),
		Client:    client,
	}, nil
}
