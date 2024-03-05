package canned

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
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
		//Image: getEnvString("ES_CONTAINER_IMAGE", "docker.elastic.co/elasticsearch/elasticsearch-oss:7.0.0"),
		Image: getEnvString("ES_CONTAINER_IMAGE", "157529275398.dkr.ecr.ap-south-1.amazonaws.com/ci-libraries/docker.elastic.co/elasticsearch/elasticsearch:6.4.2"),
		Env: map[string]string{
			"discovery.type":    "single-node",
			"network.host":      "0.0.0.0",
			"network.bind_host": "0.0.0.0",
		},
		ExposedPorts: []string{"9200/tcp", "9300/tcp"},
		//WaitingFor:   wait.ForListeningPort("9200").WithStartupTimeout(time.Minute * 3), // Default timeout is 1 minute
		WaitingFor:   wait.ForHTTP("/").WithPort("9200/tcp").WithStartupTimeout(time.Minute * 3),
		RegistryCred: getBasicAuth(),
		AutoRemove:   true,
		SkipReaper:   skipReaper(),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		ProviderType:     testContainerProvider(),
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
