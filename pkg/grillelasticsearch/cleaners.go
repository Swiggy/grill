package grillelasticsearch

import (
	"context"
	"github.com/Swiggy/grill"
	"github.com/elastic/go-elasticsearch/v7/esapi"
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

func (ge *ElasticSearch) DeleteScript(name string) grill.Cleaner {
	return grill.CleanerFunc(func() error {
		req := esapi.DeleteScriptRequest{
			ScriptID: name,
		}
		_, err := req.Do(context.Background(), ge.elasticSearch.Client)
		return err
	})
}
