package canned

import (
	"context"
	"fmt"
	"strconv"

	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
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
	env := map[string]string{
		"KAFKA_LISTENERS":                        fmt.Sprintf("PLAINTEXT://0.0.0.0:%v,BROKER://0.0.0.0:%s", kafkaPort.Port(), brokerPort.Port()),
		"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP":   "BROKER:PLAINTEXT,PLAINTEXT:PLAINTEXT",
		"KAFKA_INTER_BROKER_LISTENER_NAME":       "BROKER",
		"KAFKA_BROKER_ID":                        "1",
		"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR": "1",
		"KAFKA_OFFSETS_TOPIC_NUM_PARTITIONS":     "1",
		"KAFKA_LOG_FLUSH_INTERVAL_MESSAGES":      "10000000",
		"KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS": "0",
		"KAFKA_AUTO_CREATE_TOPICS_ENABLE":        strconv.FormatBool(true),
	}

	req := testcontainers.ContainerRequest{
		Image:        "confluentinc/cp-kafka:5.2.1",
		ExposedPorts: []string{brokerPort.Port(), kafkaPort.Port(), zookeeperPort},
		Cmd:          []string{"sleep", "infinity"},
		Env:          env,
		AutoRemove:   true,
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	err = startZookeeper(ctx, container)
	if err != nil {
		return nil, fmt.Errorf("failed to start zookeeper, error: %v", err)
	}

	port, _ := container.MappedPort(ctx, kafkaPort)
	host, _ := container.Host(ctx)
	bootstrapServer := fmt.Sprintf("PLAINTEXT://%s:%s", host, port.Port())

	err = startKafka(ctx, container, bootstrapServer)
	if err != nil {
		return nil, fmt.Errorf("failed to start kafka, error: %v", err)
	}

	ac, err := confluent.NewAdminClient(&confluent.ConfigMap{"bootstrap.servers": bootstrapServer})
	if err != nil {
		return nil, fmt.Errorf("error creating admin client, error: %v", err)
	}

	producer, err := confluent.NewProducer(&confluent.ConfigMap{"bootstrap.servers": bootstrapServer})
	if err != nil {
		return nil, fmt.Errorf("error creating producer, error: %v", err)
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("error creating docker client, error: %v", err)
	}

	return &Kafka{
		Container:    container,
		DockerClient: dockerClient,

		Host:             host,
		Port:             port.Port(),
		BootstrapServers: bootstrapServer,
		AdminClient:      ac,
		Producer:         producer,
	}, nil
}

func startZookeeper(ctx context.Context, container testcontainers.Container) error {
	cmd := []string{"sh", "-c", "printf 'clientPort=2181\ndataDir=/var/lib/zookeeper/data\ndataLogDir=/var/lib/zookeeper/log'> /zookeeper.properties; zookeeper-server-start /zookeeper.properties >/dev/null 2>&1 &"}
	_, err := container.Exec(ctx, cmd)
	if err != nil {
		return err
	}
	return nil
}

func startKafka(ctx context.Context, container testcontainers.Container, bootstrapServers string) error {
	zookeeperConnect := fmt.Sprintf("localhost:%s", zookeeperPort)
	cmd := []string{"sh", "-c", fmt.Sprintf("export KAFKA_ZOOKEEPER_CONNECT=%s; export KAFKA_ADVERTISED_LISTENERS=%s,BROKER://:%s; /etc/confluent/docker/run >/tmp/kafka-start.log 2>/tmp/kafka-start &", zookeeperConnect, bootstrapServers, brokerPort.Port())}
	_, err := container.Exec(ctx, cmd)
	if err != nil {
		return err
	}
	return nil
}
