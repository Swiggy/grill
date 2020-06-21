package grillhttp

import (
	"fmt"
	"net/http"

	"bitbucket.org/swigy/grill"
)

func (grillhttp *GrillHTTP) ResetAllStubs() grill.Cleaner {
	return grill.CleanerFunc(func() error {
		url := fmt.Sprintf("%s/__admin/mappings/reset", grillhttp.wiremock.AdminEndpoint)
		res, err := http.Post(url, "application/json", nil)
		if res != nil && res.Body != nil {
			defer res.Body.Close()
		}
		return err
	})
}
