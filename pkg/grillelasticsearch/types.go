package grillelasticsearch

import (
	"encoding/json"
)

const (
	itemFoundStatusCode    = 200
	itemNotFoundStatusCode = 404
)

type getItem struct {
	Source json.RawMessage `json:"_source"`
}
