package canned

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"time"
)

type ElasticSearch struct {
	Container testcontainers.Container
	Host      string
	Port      string
	Endpoint  string
	Client    *elasticsearch.Client
}

func NewElasticSearch(ctx context.Context) (*ElasticSearch, error) {
	_ = os.Setenv("TC_HOST", "localhost")
	req := testcontainers.ContainerRequest{
		Image: getEnvString("ES_CONTAINER_IMAGE", "elasticsearch:7.17.9"),
		Env: map[string]string{
			"discovery.type": "single-node"},
		ExposedPorts: []string{"9200/tcp"},
		WaitingFor:   wait.ForListeningPort("9200").WithStartupTimeout(time.Minute * 3), // Default timeout is 1 minute
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
	port, _ := container.MappedPort(ctx, "9200")
	endpoint, err := container.Endpoint(ctx, "")
	endpoint = fmt.Sprintf("http://%s", endpoint)

	if err != nil {
		return nil, err
	}

	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{endpoint},
	})
	if err != nil {
		return nil, err
	}

	es := &ElasticSearch{
		Container: container,
		Host:      host,
		Port:      port.Port(),
		Endpoint:  endpoint,
		Client:    client,
	}
	return es, nil
}
