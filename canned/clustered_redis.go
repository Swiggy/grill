package canned

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/go-redis/redis/v8"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"strconv"
	"time"
)

type ClusteredRedis struct {
	Host         string
	Port         string
	RedisClient  redis.UniversalClient
	Container    testcontainers.Container
	DockerClient *client.Client
}

func NewClusteredRedis() (*ClusteredRedis, error) {
	os.Setenv("TC_HOST", "localhost")
	totalNumberOfRedisClusters := 3
	allAvailablePorts, err := getConsecutiveFreePorts(totalNumberOfRedisClusters, 49152, 55000)
	if err != nil {
		fmt.Errorf("failed to start Clustered Redis, error: %v", err)
		return nil, err
	}
	startingPort := allAvailablePorts[0]
	endPort := allAvailablePorts[totalNumberOfRedisClusters-1]
	startingPortStr := strconv.Itoa(startingPort)
	exposedPorts := fmt.Sprintf("%v-%v:%v-%v", startingPort, endPort, startingPort, endPort)

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        getEnvString("CLUSTERED_REDIS_CONTAINER_IMAGE", "grokzen/redis-cluster:6.2.1"),
		ExposedPorts: []string{exposedPorts},
		WaitingFor: wait.ForLog("Cluster state changed: ok").
			WithOccurrence(3).
			WithPollInterval(time.Second * 5),
		Env: map[string]string{
			"INITIAL_PORT":      startingPortStr,
			"REDIS_CLUSTER_IP":  "0.0.0.0",
			"IP":                "127.0.0.1",
			"MASTERS":           strconv.Itoa(totalNumberOfRedisClusters),
			"SLAVES_PER_MASTER": "0",
		},
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
	port, _ := container.MappedPort(ctx, nat.Port(startingPortStr))

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("error creating docker client for RedisCluster, error: %v", err)
	}

	redisClient, err := initializeRedisClient(host, startingPortStr)
	if err != nil {
		return nil, fmt.Errorf("error creating Clustered RedisClient, error: %v", err)
	}
	return &ClusteredRedis{
		Host:         host,
		Port:         port.Port(),
		Container:    container,
		DockerClient: dockerClient,
		RedisClient:  redisClient,
	}, nil
}

func initializeRedisClient(host string, port string) (redis.UniversalClient, error) {
	redisAddressURL := fmt.Sprintf("%s:%s", host, port)

	redisClusterOption := redis.UniversalOptions{
		Addrs: []string{
			redisAddressURL,
		},
		ReadTimeout:  500 * time.Millisecond,
		WriteTimeout: 500 * time.Millisecond,
		PoolSize:     5,
		MaxRetries:   3,
		MinIdleConns: 3,
	}
	newRedisClient := redis.NewClusterClient(redisClusterOption.Cluster())
	return newRedisClient, nil
}
