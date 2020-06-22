package grilldp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"bitbucket.org/swigy/grill/internal/canned"
)

type GrillDP struct {
	wiremock *canned.WireMock
}

func Start() (*GrillDP, error) {
	wiremock, err := canned.NewWiremock(context.TODO())
	if err != nil {
		return nil, err
	}

	for _, stub := range []stub{registerEventStub, messageSetStub, messageStub} {
		if err := createStub(stub, wiremock.AdminEndpoint); err != nil {
			return nil, fmt.Errorf("error creating stub, error=%v", err)
		}
	}

	return &GrillDP{
		wiremock: wiremock,
	}, nil
}

func (gdp *GrillDP) Host() string {
	return gdp.wiremock.Host
}

func (gdp *GrillDP) Port() string {
	return gdp.wiremock.Port
}

func (gdp *GrillDP) Stop() error {
	return gdp.wiremock.Container.Terminate(context.TODO())
}

func createStub(stub stub, adminEndpoint string) error {
	url := fmt.Sprintf("%s/__admin/mappings", adminEndpoint)
	jsonStr, err := json.Marshal(stub)
	if err != nil {
		return err
	}
	res, err := http.Post(url, "application/json", strings.NewReader(string(jsonStr)))
	if res != nil && res.Body != nil {
		defer res.Body.Close()
	}
	return err
}
