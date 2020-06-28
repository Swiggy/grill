package canned

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/docker/docker/client"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/net/http2"
)

type DynamoDB struct {
	Container    testcontainers.Container
	DockerClient *client.Client

	Host      string
	Port      string
	AccessKey string
	SecretKey string
	Region    string
	Client    dynamodbiface.DynamoDBAPI
}

func NewDynamoDB(ctx context.Context) (*DynamoDB, error) {
	os.Setenv("TC_HOST", "localhost")
	req := testcontainers.ContainerRequest{
		Image:        "amazon/dynamodb-local",
		ExposedPorts: []string{"8000/tcp"},
		WaitingFor:   wait.ForListeningPort("8000"),
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
	port, _ := container.MappedPort(ctx, "8000")
	accessKey, secretKey, region := "random", "random", "ap-southeast-1"
	endpoint := fmt.Sprintf("http://%s:%s", host, port.Port())

	dynamoClient, err := newDynamoClient(endpoint, accessKey, secretKey, region)
	if err != nil {
		return nil, err
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("error creating docker client, error: %v", err)
	}

	return &DynamoDB{
		Container:    container,
		DockerClient: dockerClient,

		Host:      host,
		Port:      port.Port(),
		AccessKey: accessKey,
		SecretKey: secretKey,
		Region:    region,
		Client:    dynamoClient,
	}, nil
}

func newDynamoClient(endpoint, accessKey, secretKey, region string) (dynamodbiface.DynamoDBAPI, error) {
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

	mySession, err := session.NewSession(&awsConfig)
	if err != nil {
		return nil, err
	}
	var svc = dynamodb.New(mySession)
	return dynamodbiface.DynamoDBAPI(svc), nil
}
