package grillhttp

import (
	"fmt"
	"net/http"

	"github.com/swiggy-private/grill"
)

func (gh *HTTP) ResetAllStubs() grill.Cleaner {
	return grill.CleanerFunc(func() error {
		url := fmt.Sprintf("%s/__admin/mappings/reset", gh.wiremock.AdminEndpoint)
		res, err := http.Post(url, "application/json", nil)
		if res != nil && res.Body != nil {
			defer res.Body.Close()
		}
		return err
	})
}
