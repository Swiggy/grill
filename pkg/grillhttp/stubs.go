package grillhttp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/lovlin-thakkar/swiggy-grill"
)

func (gh *HTTP) StubFromJSON(stubStr string) grill.Stub {
	return grill.StubFunc(func() error {
		url := fmt.Sprintf("%s/__admin/mappings", gh.wiremock.AdminEndpoint)
		res, err := http.Post(url, "application/json", strings.NewReader(stubStr))
		if res != nil && res.Body != nil {
			defer res.Body.Close()
		}
		return err
	})
}

func (gh *HTTP) Stub(request Request, response Response) grill.Stub {
	return grill.StubFunc(func() error {
		jsonStr, err := json.Marshal(&Stub{Request: request, Response: response})
		if err != nil {
			return err
		}
		return gh.StubFromJSON(string(jsonStr)).Stub()
	})
}

func (gh *HTTP) StubFromFile(filepath string) grill.Stub {
	return grill.StubFunc(func() error {
		stubStr, err := ioutil.ReadFile(filepath)
		if err != nil {
			return err
		}
		return gh.StubFromJSON(string(stubStr)).Stub()
	})
}
