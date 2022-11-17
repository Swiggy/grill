package canned

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/client"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type WireMock struct {
	Container    testcontainers.Container
	DockerClient *client.Client

	Host          string
	Port          string
	AdminEndpoint string
}

func NewWiremock(ctx context.Context) (*WireMock, error) {
	os.Setenv("TC_HOST", "localhost")
	req := testcontainers.ContainerRequest{
		Image:        "rodolpheche/wiremock",
		SkipReaper:   true,
		ExposedPorts: []string{"8080/tcp", "8443/tcp"},
		WaitingFor:   wait.ForListeningPort("8080"),
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
	port, _ := container.MappedPort(ctx, "8080")

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("error creating docker client, error: %v", err)
	}

	return &WireMock{
		Container:    container,
		DockerClient: dockerClient,

		Host:          host,
		Port:          port.Port(),
		AdminEndpoint: fmt.Sprintf("http://%s:%s", host, port.Port()),
	}, nil
}
