package canned

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Consul struct {
	Container testcontainers.Container

	Host   string
	Port   string
	Client *api.Client
}

func NewConsul(ctx context.Context) (*Consul, error) {
	req := testcontainers.ContainerRequest{
		Image:        "consul:1.7.3",
		ExposedPorts: []string{"8500/tcp"},
		WaitingFor:   wait.ForListeningPort("8500"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, err
	}

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "8500")

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", host, port.Port())
	client, err := api.NewClient(config)

	if err != nil {
		return nil, fmt.Errorf("error creating client")
	}

	return &Consul{
		Host:      host,
		Port:      port.Port(),
		Container: container,
		Client:    client,
	}, nil
}
