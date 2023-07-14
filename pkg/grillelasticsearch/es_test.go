package grillelasticsearch

import (
	"context"
	"encoding/json"
	"github.com/Swiggy/grill"
	"testing"
)

func TestElasticSearch_PutItem(t *testing.T) {
	helper := &ElasticSearch{}
	if err := helper.Start(context.TODO()); err != nil {
		t.Errorf("error starting elastic search grill, error=%v", err)
		return
	}
	const testData = `{"fname": "John", "lname": "Doe"}`

	tests := []grill.TestCase{
		{
			Name: "Item should get indexed",
			Stubs: []grill.Stub{
				helper.CreateIndex("test_index"),
				helper.PutItem("test_index", "doc_id", testData),
			},
			Action: func() interface{} {
				return nil
			},
			Assertions: []grill.Assertion{
				helper.AssertItem("test_index", "doc_id", json.RawMessage(testData)),
			},
			Cleaners: []grill.Cleaner{
				helper.DeleteIndices("test_index"),
			},
		},
		{
			Name: "Item not found",
			Stubs: []grill.Stub{
				helper.CreateIndex("test_index"),
			},
			Action: func() interface{} {
				return nil
			},
			Assertions: []grill.Assertion{
				helper.AssertNoItem("test_index", "doc_id"),
			},
			Cleaners: []grill.Cleaner{
				helper.DeleteIndices("test_index"),
			},
		},
	}

	grill.Run(t, tests)
}
