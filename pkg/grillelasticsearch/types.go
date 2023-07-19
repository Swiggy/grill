package grillelasticsearch

import (
	"encoding/json"
)

const (
	itemFoundStatusCode          = 200
	itemNotFoundStatusCode       = 404
	resourceModifiedSuccessfully = 200
	resourceCreatedSuccessfully  = 201
)

type getItem struct {
	Source json.RawMessage `json:"_source"`
}
