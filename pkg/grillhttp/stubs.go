package grillhttp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"bitbucket.org/swigy/grill"
)

func (grillhttp *GrillHTTP) StubFromJSON(stubStr string) grill.Stub {
	return grill.StubFunc(func() error {
		url := fmt.Sprintf("%s/__admin/mappings", grillhttp.wiremock.AdminEndpoint)
		res, err := http.Post(url, "application/json", strings.NewReader(stubStr))
		if res != nil && res.Body != nil {
			defer res.Body.Close()
		}
		return err
	})
}

func (grillhttp *GrillHTTP) Stub(stub *Stub) grill.Stub {
	return grill.StubFunc(func() error {
		jsonStr, err := json.Marshal(stub)
		if err != nil {
			return err
		}
		return grillhttp.StubFromJSON(string(jsonStr)).Stub()
	})
}

func (grillhttp *GrillHTTP) StubFromFile(filepath string) grill.Stub {
	return grill.StubFunc(func() error {
		stubStr, err := ioutil.ReadFile(filepath)
		if err != nil {
			return err
		}
		return grillhttp.StubFromJSON(string(stubStr)).Stub()
	})
}
