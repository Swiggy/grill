package grillelasticsearch

import (
	"encoding/json"
)

const (
	resourceModifiedSuccessfully = 200
	resourceCreatedSuccessfully  = 201
)

type searchItems struct {
	Hits *searchHits `json:"hits,omitempty"` // the actual search hits
}

// searchHits specifies the list of search hits.
type searchHits struct {
	Hits       []*searchHitSource `json:"hits,omitempty"` // the actual hits returned
	TotalCount totalCount         `json:"total,omitempty"`
}

type totalCount struct {
	Value int `json:"value,omitempty"`
}

// searchHitSource is a single hit
type searchHitSource struct {
	Source json.RawMessage `json:"_source"`
}
