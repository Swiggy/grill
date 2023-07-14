package grillelasticsearch

import (
	"github.com/Swiggy/grill"
	"strings"
)

func (ge *ElasticSearch) CreateIndex(index string) grill.Stub {
	return grill.StubFunc(func() error {
		_, err := ge.elasticSearch.Client.Indices.Create(index)
		return err
	})
}

func (ge *ElasticSearch) PutItem(index, docId, data string) grill.Stub {
	return grill.StubFunc(func() error {
		_, err := ge.elasticSearch.Client.Create(index, docId, strings.NewReader(data))
		return err
	})
}
