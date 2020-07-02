package grilldp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"bitbucket.org/swigy/grill"
	"github.com/xeipuuv/gojsonschema"
)

func (gdp *DP) AssertRegisteredApps(apps ...string) grill.Assertion {
	return grill.AssertionFunc(func() error {
		journal, err := gdp.getMatchingRequests(registerEventStub.Request)
		if err != nil || journal == nil {
			return fmt.Errorf("no registered apps found")
		}
		var result []string
		for _, match := range journal.Requests {
			event := uuidRegister{}
			json.Unmarshal([]byte(match.Body), &event)
			result = append(result, event.AppName)
		}

		if len(result) != len(apps) {
			return fmt.Errorf("invalid apps registed, got=%v, want=%v", result, apps)
		}

		sort.Slice(result, func(i, j int) bool {
			return result[i] < result[j]
		})
		sort.Slice(apps, func(i, j int) bool {
			return apps[i] < apps[j]
		})

		for i := 0; i < len(result); i++ {
			if result[i] != apps[i] {
				return fmt.Errorf("invalid apps registed, got=%v, want=%v", result, apps)
			}
		}

		return nil
	})

}

func (gdp *DP) AssertCount(appName, eventName, version string, expectedCount int) grill.Assertion {
	return grill.AssertionFunc(func() error {
		journal, err := gdp.getMatchingRequests(messageSetStub.Request)
		if err != nil || journal == nil {
			return nil
		}
		count := 0
		for _, request := range journal.Requests {
			var events []*event
			json.Unmarshal([]byte(request.Body), &events)
			for _, event := range events {
				if event.Header.AppName == appName && event.Header.Name == eventName && event.Header.SchemaVersion == version {
					count++
				}
			}
		}

		if count != expectedCount {
			return fmt.Errorf("invalid count for event=%v version=%v, got=%v, want=%v", eventName, version, count, expectedCount)
		}

		return nil
	})
}

func (gdp *DP) AssertValidSchema(appName, eventName, version, schemaFilePath string) grill.Assertion {
	return grill.AssertionFunc(func() error {
		schemaLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", schemaFilePath))
		schema, err := gojsonschema.NewSchema(schemaLoader)
		if err != nil {
			return fmt.Errorf("error loading schema at path=%v, error=%v", schemaFilePath, err)
		}

		journal, err := gdp.getMatchingRequests(messageSetStub.Request)
		if err != nil || journal == nil {
			return nil
		}

		var failures []string
		for _, request := range journal.Requests {
			var events []*event
			json.Unmarshal([]byte(request.Body), &events)
			for _, event := range events {
				if event.Header.AppName == appName && event.Header.Name == eventName && event.Header.SchemaVersion == version {
					eventLoader := gojsonschema.NewGoLoader(event.Event)
					if res, err := schema.Validate(eventLoader); err != nil || !res.Valid() {
						for _, e := range res.Errors() {
							failures = append(failures, "error: %v", e.String())
						}
						if err != nil {
							failures = append(failures, "error: %v", err.Error())
						}
					}
				}
			}
		}

		if len(failures) > 0 {
			return fmt.Errorf("%v events failed validation, failures=%v", len(failures), failures)
		}

		return nil
	})
}

func (gdp *DP) getMatchingRequests(req request) (*requestJournal, error) {
	url := fmt.Sprintf("%s/__admin/requests/find", gdp.wiremock.AdminEndpoint)
	jsonStr, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	res, err := http.Post(url, "application/json", strings.NewReader(string(jsonStr)))
	if err != nil {
		return nil, err
	}

	if res == nil || res.Body == nil {
		return nil, nil
	}
	defer res.Body.Close()

	jj := requestJournal{}
	json.NewDecoder(res.Body).Decode(&jj)

	if len(jj.Requests) == 0 {
		return nil, nil
	}

	return &jj, nil
}
