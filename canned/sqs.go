package canned

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/docker/docker/client"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/net/http2"
	"net"
	"net/http"
	"os"
)

type SQS struct {
	Container    testcontainers.Container
	DockerClient *client.Client

	Host      string
	Port      string
	AccessKey string
	SecretKey string
	Region    string
	Client    sqsiface.SQSAPI
}

func NewSQS(ctx context.Context) (*SQS, error) {
	os.Setenv("TC_HOST", "localhost")
	req := testcontainers.ContainerRequest{
		Image:        getEnvString("SQS_CONTAINER_IMAGE", "softwaremill/elasticmq-native:1.3.14"),
		ExposedPorts: []string{"9324/tcp"},
		WaitingFor:   wait.ForListeningPort("9324"),
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
	port, _ := container.MappedPort(ctx, "9324")
	accessKey, secretKey, region := getAWSConfig()
	endpoint := fmt.Sprintf("http://%s:%s", host, port.Port())

	sqsClient, err := newSQSClient(endpoint, accessKey, secretKey, region)
	if err != nil {
		return nil, err
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("error creating docker client, error: %v", err)
	}

	return &SQS{
		Container:    container,
		DockerClient: dockerClient,

		Host:      host,
		Port:      port.Port(),
		AccessKey: accessKey,
		SecretKey: secretKey,
		Region:    region,
		Client:    sqsClient,
	}, nil
}

func newSQSClient(endpoint, accessKey, secretKey, region string) (sqsiface.SQSAPI, error) {
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
	var svc = sqs.New(mySession)
	return sqsiface.SQSAPI(svc), nil
}
