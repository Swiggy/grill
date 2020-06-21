package grillconsul

import (
	"context"

	"bitbucket.org/swigy/grill/internal/canned"
	"github.com/hashicorp/consul/api"
)

type GrillConsul struct {
	consul *canned.Consul
}

func Start() (*GrillConsul, error) {
	consul, err := canned.NewConsul(context.TODO())
	if err != nil {
		return nil, err
	}

	return &GrillConsul{
		consul: consul,
	}, nil
}

func (gc *GrillConsul) Host() string {
	return gc.consul.Host
}

func (gc *GrillConsul) Port() string {
	return gc.consul.Port
}

func (gc *GrillConsul) Client() *api.Client {
	return gc.consul.Client
}

func (gc *GrillConsul) Stop() error {
	return gc.consul.Container.Terminate(context.Background())
}
