package canned

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/docker/docker/client"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Minio struct {
	Container    testcontainers.Container
	DockerClient *client.Client

	Host      string
	Port      string
	AccessKey string
	SecretKey string
	Region    string
	Client    s3iface.S3API
}

func NewMinio(ctx context.Context) (*Minio, error) {
	os.Setenv("TC_HOST", "localhost")
	req := testcontainers.ContainerRequest{
		Image:        "minio/minio",
		ExposedPorts: []string{"9000/tcp"},
		WaitingFor:   wait.ForListeningPort("9000/tcp"),
		Cmd:          []string{"server", "/data"},
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
	port, _ := container.MappedPort(ctx, "9000")
	accessKey, secretKey, region := "random", "random", "ap-southeast-1"

	s3Endpoint := fmt.Sprintf("http://%s:%s", host, port.Port())
	awsSession, err := session.NewSession(&aws.Config{
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(endpoints.ApSoutheast1RegionID),
		Endpoint:         aws.String(s3Endpoint),
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		return nil, err
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("error creating docker client, error: %v", err)
	}

	return &Minio{
		Container:    container,
		DockerClient: dockerClient,

		Host:      host,
		Port:      port.Port(),
		AccessKey: accessKey,
		SecretKey: secretKey,
		Region:    region,
		Client:    s3.New(awsSession),
	}, nil
}
