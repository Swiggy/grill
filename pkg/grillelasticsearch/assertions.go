package grillelasticsearch

import (
	"encoding/json"
	"fmt"
	"github.com/Swiggy/grill"
	"io"
	"reflect"
)

func (ge *ElasticSearch) AssertItem(index, docId string, expected json.RawMessage) grill.Assertion {
	return grill.AssertionFunc(func() error {
		output, err := ge.elasticSearch.Client.Get(index, docId)
		if err != nil {
			return err
		}

		if output.StatusCode != itemFoundStatusCode {
			return fmt.Errorf("invalid status code, got=%v, want=%v ", output.StatusCode, itemFoundStatusCode)
		}

		bytes, err := io.ReadAll(output.Body)
		if err != nil {
			return err
		}

		var item getItem
		err = json.Unmarshal(bytes, &item)
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(item.Source, expected) {
			return fmt.Errorf("invalid item, got=%v, want=%v", item.Source, expected)
		}
		return nil
	})
}

func (ge *ElasticSearch) AssertNoItem(index, docId string) grill.Assertion {
	return grill.AssertionFunc(func() error {
		output, err := ge.elasticSearch.Client.Get(index, docId)
		if err != nil {
			return err
		}

		if output.StatusCode != itemNotFoundStatusCode {
			return fmt.Errorf("item exits")
		}

		return nil
	})
}
