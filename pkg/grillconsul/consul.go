package grillconsul

import (
	"context"

	"github.com/swiggy-private/grill/internal/canned"
	"github.com/hashicorp/consul/api"
)

type Consul struct {
	consul *canned.Consul
}

func (gc *Consul) Start(ctx context.Context) error {
	consul, err := canned.NewConsul(ctx)
	if err != nil {
		return err
	}

	gc.consul = consul
	return nil
}

func (gc *Consul) Host() string {
	return gc.consul.Host
}

func (gc *Consul) Port() string {
	return gc.consul.Port
}

func (gc *Consul) Client() *api.Client {
	return gc.consul.Client
}

func (gc *Consul) Stop(ctx context.Context) error {
	return gc.consul.Container.Terminate(ctx)
}
