package grillhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/lovlin-thakkar/swiggy-grill"
)

func (gh *HTTP) AssertCount(request *Request, expectedCount int) grill.Assertion {
	return grill.AssertionFunc(func() error {
		url := fmt.Sprintf("%s/__admin/requests/count", gh.wiremock.AdminEndpoint)
		jsonStr, err := json.Marshal(request)
		if err != nil {
			return err
		}
		res, err := http.Post(url, "application/json", strings.NewReader(string(jsonStr)))
		if res == nil && res.Body == nil {
			return nil
		}

		var output struct {
			Count int `json:"count"`
		}
		err = json.NewDecoder(res.Body).Decode(&output)
		if err != nil {
			return err
		}

		if output.Count != expectedCount {
			return fmt.Errorf("invalid number of requests, got=%v, want=%v", output.Count, expectedCount)
		}

		return nil
	})
}
