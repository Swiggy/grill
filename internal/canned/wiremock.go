package canned

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type WireMock struct {
	Container testcontainers.Container

	Host          string
	Port          string
	AdminEndpoint string
}

func NewWiremock(ctx context.Context) (*WireMock, error) {
	req := testcontainers.ContainerRequest{
		Image:        "rodolpheche/wiremock",
		ExposedPorts: []string{"8080/tcp", "8443/tcp"},
		WaitingFor:   wait.ForListeningPort("8080"),
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

	return &WireMock{
		Host:          host,
		Port:          port.Port(),
		Container:     container,
		AdminEndpoint: fmt.Sprintf("http://%s:%s", host, port.Port()),
	}, nil
}
