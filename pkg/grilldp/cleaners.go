package grilldp

import (
	"fmt"
	"net/http"

	"bitbucket.org/swigy/grill"
)

func (gdp *DP) FlushAllEvents() grill.Cleaner {
	return grill.CleanerFunc(func() error {
		url := fmt.Sprintf("%s/__admin/requests", gdp.wiremock.AdminEndpoint)
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		res, err := (&http.Client{}).Do(req)
		if res != nil && res.Body != nil {
			defer res.Body.Close()
		}
		return err
	})
}
