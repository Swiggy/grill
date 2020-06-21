package grillhttp

import (
	"context"

	"bitbucket.org/swigy/grill/internal/canned"
)

type GrillHTTP struct {
	wiremock *canned.WireMock
}

func Start() (*GrillHTTP, error) {
	wiremock, err := canned.NewWiremock(context.TODO())
	if err != nil {
		return nil, err
	}
	return &GrillHTTP{
		wiremock: wiremock,
	}, nil
}

func (grillhttp *GrillHTTP) Host() string {
	return grillhttp.wiremock.Host
}

func (grillhttp *GrillHTTP) Port() string {
	return grillhttp.wiremock.Port
}

func (grillhttp *GrillHTTP) AdminEndpoint() string {
	return grillhttp.wiremock.AdminEndpoint
}

func (grillhttp *GrillHTTP) Stop() error {
	return grillhttp.wiremock.Container.Terminate(context.TODO())
}
