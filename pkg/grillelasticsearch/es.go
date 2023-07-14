package grillelasticsearch

import (
	"context"
	"github.com/Swiggy/grill/canned"
	"github.com/elastic/go-elasticsearch/v7"
)

type ElasticSearch struct {
	elasticSearch *canned.ElasticSearch
}

func (ge *ElasticSearch) Start(ctx context.Context) error {
	elasticSearch, err := canned.NewElasticSearch(ctx)
	if err != nil {
		return err
	}
	ge.elasticSearch = elasticSearch
	return nil
}

func (ge *ElasticSearch) Client() *elasticsearch.Client {
	return ge.elasticSearch.Client
}

func (ge *ElasticSearch) Host() string {
	return ge.elasticSearch.Host
}

func (ge *ElasticSearch) Port() string {
	return ge.elasticSearch.Port
}

func (ge *ElasticSearch) Stop(ctx context.Context) error {
	return ge.elasticSearch.Container.Terminate(ctx)
}
