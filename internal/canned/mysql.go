package canned

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Mysql struct {
	Container testcontainers.Container

	Host     string
	Port     string
	Database string
	Username string
	Password string
	Client   *sql.DB
}

func NewMysql(ctx context.Context) (*Mysql, error) {
	req := testcontainers.ContainerRequest{
		Image:        "mysql:5.6",
		ExposedPorts: []string{"3306/tcp"},
		WaitingFor:   wait.ForListeningPort("3306"),
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "password",
			"MYSQL_DATABASE":      "test",
		},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "3306")

	client, err := sql.Open("mysql", fmt.Sprintf("root:password@tcp(%s:%s)/test", host, port.Port()))
	if err != nil {
		return nil, err
	}

	if _, err := client.Exec("USE DATABASE test"); err != nil {
		return nil, fmt.Errorf("error setting database test, error: %v", err)
	}

	return &Mysql{
		Container: container,
		Host:      host,
		Port:      port.Port(),
		Database:  "test",
		Username:  "root",
		Password:  "password",
		Client:    client,
	}, nil
}
