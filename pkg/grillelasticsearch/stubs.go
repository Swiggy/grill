package grillelasticsearch

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Swiggy/grill"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"strings"
)

func (ge *ElasticSearch) CreateIndex(index string, mapping string) grill.Stub {
	return grill.StubFunc(func() error {
		req := esapi.IndicesCreateRequest{
			Index: index,
			Body:  strings.NewReader(mapping),
		}
		res, err := req.Do(context.Background(), ge.elasticSearch.Client)
		if err != nil {
			return err
		}
		if res.StatusCode != resourceModifiedSuccessfully {
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(res.Body)
			respBytes := buf.String()
			return fmt.Errorf(respBytes)
		}
		return nil
	})
}

func (ge *ElasticSearch) UpsertItem(index, docId, data string) grill.Stub {
	return grill.StubFunc(func() error {
		req := esapi.IndexRequest{
			Index:      index,
			Body:       strings.NewReader(data),
			DocumentID: docId,
		}
		res, err := req.Do(context.Background(), ge.elasticSearch.Client)
		if err != nil {
			return err
		}
		if res.StatusCode != resourceModifiedSuccessfully && res.StatusCode != resourceCreatedSuccessfully {
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(res.Body)
			respBytes := buf.String()
			return fmt.Errorf(respBytes)
		}
		return nil
	})
}

func (ge *ElasticSearch) AddScript(name string, script string) grill.Stub {
	return grill.StubFunc(func() error {
		req := esapi.PutScriptRequest{
			Body:     strings.NewReader(script),
			ScriptID: name,
		}
		res, err := req.Do(context.Background(), ge.elasticSearch.Client)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		if res.StatusCode != resourceModifiedSuccessfully && res.StatusCode != resourceCreatedSuccessfully {
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(res.Body)
			respBytes := buf.String()
			return fmt.Errorf(respBytes)
		}
		return nil
	})
}
