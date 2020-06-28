package grillhttp

import (
	"context"

	"bitbucket.org/swigy/grill/internal/canned"
)

type HTTP struct {
	wiremock *canned.WireMock
}

func (gh *HTTP) Start(ctx context.Context) error {
	wiremock, err := canned.NewWiremock(ctx)
	if err != nil {
		return err
	}
	gh.wiremock = wiremock
	return nil
}

func (gh *HTTP) Host() string {
	return gh.wiremock.Host
}

func (gh *HTTP) Port() string {
	return gh.wiremock.Port
}

func (gh *HTTP) AdminEndpoint() string {
	return gh.wiremock.AdminEndpoint
}

func (gh *HTTP) Stop() error {
	return gh.wiremock.Container.Terminate(context.TODO())
}
