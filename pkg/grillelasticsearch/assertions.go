package grillelasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Swiggy/grill"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"io"
	"reflect"
)

func (ge *ElasticSearch) AssertItemsCount(index string, count int, expected []json.RawMessage) grill.Assertion {
	return grill.AssertionFunc(func() error {

		req := esapi.SearchRequest{
			Index: []string{index},
			Body:  nil,
		}
		resp, err := req.Do(context.Background(), ge.elasticSearch.Client)
		if err != nil {
			return err
		}
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var items searchItems
		err = json.Unmarshal(bytes, &items)
		if err != nil {
			return err
		}

		if len(items.Hits.Hits) != count {
			return fmt.Errorf("unequal count, got=%v, want=%v ", len(items.Hits.Hits), count)
		}

		if count == 0 { // no need to check items
			return nil
		}

		rawItems := make([]json.RawMessage, 0)
		for _, item := range items.Hits.Hits {
			rawItems = append(rawItems, item.Source)
		}

		if !reflect.DeepEqual(rawItems, expected) {
			return fmt.Errorf("invalid item, got=%v, want=%v", rawItems, expected)
		}
		return nil
	})
}
