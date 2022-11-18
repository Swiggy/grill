package canned

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/docker/docker/client"
	"github.com/hashicorp/consul/api"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Consul struct {
	Container    testcontainers.Container
	DockerClient *client.Client

	Host   string
	Port   string
	Client *api.Client
}

func NewConsul(ctx context.Context) (*Consul, error) {
	os.Setenv("TC_HOST", "localhost")
	skipReaper, _ := strconv.ParseBool(os.Getenv("TESTCONTAINERS_RYUK_DISABLED"))
	req := testcontainers.ContainerRequest{
		Image:        "consul:1.7.3",
		SkipReaper:   skipReaper,
		ExposedPorts: []string{"8500/tcp"},
		WaitingFor:   wait.ForListeningPort("8500"),
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
	port, _ := container.MappedPort(ctx, "8500")

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", host, port.Port())
	consulClient, err := api.NewClient(config)

	if err != nil {
		return nil, fmt.Errorf("error creating client")
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("error creating docker client, error: %v", err)
	}

	return &Consul{
		Container:    container,
		DockerClient: dockerClient,

		Host:   host,
		Port:   port.Port(),
		Client: consulClient,
	}, nil
}
