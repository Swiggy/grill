package grillelasticsearch

import (
	"github.com/Swiggy/grill"
)

func (ge *ElasticSearch) DeleteIndices(indices ...string) grill.Cleaner {
	return grill.CleanerFunc(func() error {
		_, err := ge.elasticSearch.Client.Indices.Delete(indices)
		return err
	})
}

func (ge *ElasticSearch) DeleteItem(index, id string) grill.Cleaner {
	return grill.CleanerFunc(func() error {
		_, err := ge.elasticSearch.Client.Delete(index, id)
		return err
	})
}
