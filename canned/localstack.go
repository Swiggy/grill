package canned

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/docker/docker/client"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/net/http2"
	"net"
	"net/http"
	"os"
)

type Localstack struct {
	Container    testcontainers.Container
	DockerClient *client.Client

	Host       string
	Port       string
	AccessKey  string
	SecretKey  string
	Region     string
	AWSSession *session.Session
}

func NewLocalstack(ctx context.Context) (*Localstack, error) {
	os.Setenv("TC_HOST", "localhost")
	req := testcontainers.ContainerRequest{
		Image:        getEnvString("LOCALSTACK_CONTAINER_IMAGE", "localstack/localstack:1.0.3"),
		ExposedPorts: []string{"4566/tcp"},
		WaitingFor:   wait.ForListeningPort("4566"),
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
	port, _ := container.MappedPort(ctx, "4566")
	accessKey, secretKey, region := getAWSConfig()
	endpoint := fmt.Sprintf("http://%s:%s", host, port.Port())

	awsSession, err := newAWSSession(accessKey, secretKey, endpoint, region)
	if err != nil {
		return nil, err
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("error creating docker client, error: %v", err)
	}

	return &Localstack{
		Container:    container,
		DockerClient: dockerClient,

		Host:       host,
		Port:       port.Port(),
		Region:     region,
		AWSSession: awsSession,
	}, nil
}

func newAWSSession(accessKey, secretKey, endpoint, region string) (*session.Session, error) {
	var awsConfig aws.Config
	awsConfig.Region = aws.String(region)
	awsConfig.Endpoint = aws.String(endpoint)
	awsConfig.Credentials = credentials.NewStaticCredentials(accessKey, secretKey, "")
	awsConfig.MaxRetries = aws.Int(3)

	tr := &http.Transport{DialContext: (&net.Dialer{}).DialContext}
	if err := http2.ConfigureTransport(tr); err != nil {
		return nil, err
	}

	awsConfig.HTTPClient = &http.Client{Transport: tr}

	return session.NewSession(&awsConfig)
}
