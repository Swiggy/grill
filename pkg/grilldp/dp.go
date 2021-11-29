package grilldp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/swiggy-private/grill/internal/canned"
)

type DP struct {
	wiremock *canned.WireMock
}

func (gdp *DP) Start(ctx context.Context) error {
	wiremock, err := canned.NewWiremock(ctx)
	if err != nil {
		return err
	}

	for _, stub := range []stub{registerEventStub, messageSetStub, messageStub} {
		if err := createStub(stub, wiremock.AdminEndpoint); err != nil {
			return fmt.Errorf("error creating stub, error=%v", err)
		}
	}
	gdp.wiremock = wiremock
	return nil
}

func (gdp *DP) Host() string {
	return gdp.wiremock.Host
}

func (gdp *DP) Port() string {
	return gdp.wiremock.Port
}

func (gdp *DP) Stop(ctx context.Context) error {
	return gdp.wiremock.Container.Terminate(ctx)
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
