package canned

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
	"os"

	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	kafkaPort     nat.Port = "9093/tcp"
	brokerPort    nat.Port = "9092/tcp"
	zookeeperPort          = "2181"
)

type Kafka struct {
	Container    testcontainers.Container
	DockerClient *client.Client

	Host             string
	Port             string
	BootstrapServers string
	AdminClient      *confluent.AdminClient
	Producer         *confluent.Producer
}

func NewKafka(ctx context.Context) (*Kafka, error) {
	os.Setenv("TC_HOST", "localhost")

	kafkaContainer, err := kafka.RunContainer(ctx, testcontainers.CustomizeRequestOption(kafkaTestContainerProvider))
	if err != nil {
		return nil, err
	}

	port, _ := kafkaContainer.MappedPort(ctx, kafkaPort)
	host, _ := kafkaContainer.Host(ctx)

	bootstrapServer := fmt.Sprintf("PLAINTEXT://%s:%s", host, port.Port())

	ac, err := confluent.NewAdminClient(&confluent.ConfigMap{"bootstrap.servers": bootstrapServer, "api.version.request": false})
	if err != nil {
		return nil, fmt.Errorf("error creating admin client, error: %v", err)
	}

	producer, err := confluent.NewProducer(&confluent.ConfigMap{"bootstrap.servers": bootstrapServer, "api.version.request": false})
	if err != nil {
		return nil, fmt.Errorf("error creating producer, error: %v", err)
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("error creating docker client, error: %v", err)
	}

	return &Kafka{
		Container:    kafkaContainer,
		DockerClient: dockerClient,

		Host:             host,
		Port:             port.Port(),
		BootstrapServers: bootstrapServer,
		AdminClient:      ac,
		Producer:         producer,
	}, nil
}

func kafkaTestContainerProvider(req *testcontainers.GenericContainerRequest) {
	req.ProviderType = testContainerProvider()
}
