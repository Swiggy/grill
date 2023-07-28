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
	const mapping = `{"mappings":{"properties":{"fname":{"type":"keyword"},"lname":{"type":"keyword"}}}}`
	const testData = `{"fname": "John", "lname": "Doe"}`
	const modifiedTestData = `{"fname": "NewJohn", "lname": "Doe"}`

	tests := []grill.TestCase{
		{
			Name: "Items should get indexed",
			Stubs: []grill.Stub{
				helper.CreateIndex("test_index", mapping),
				helper.UpsertItem("test_index", "doc_id", testData),
				helper.UpsertItem("test_index", "doc_id_2", testData),
			},
			Action: func() interface{} {
				return nil
			},
			Assertions: []grill.Assertion{
				helper.AssertItemsCount("test_index", 2, []json.RawMessage{json.RawMessage(testData), json.RawMessage(testData)}),
			},
			Cleaners: []grill.Cleaner{
				helper.DeleteIndices("test_index"),
			},
		},
		{
			Name: "Item should get modified",
			Stubs: []grill.Stub{
				helper.CreateIndex("test_index", mapping),
				helper.UpsertItem("test_index", "doc_id", testData),
				helper.UpsertItem("test_index", "doc_id", modifiedTestData),
			},
			Action: func() interface{} {
				return nil
			},
			Assertions: []grill.Assertion{
				helper.AssertItemsCount("test_index", 1, []json.RawMessage{json.RawMessage(modifiedTestData)}),
			},
			Cleaners: []grill.Cleaner{
				helper.DeleteIndices("test_index"),
			},
		},
		{
			Name: "Item not found",
			Stubs: []grill.Stub{
				helper.CreateIndex("test_index", mapping),
			},
			Action: func() interface{} {
				return nil
			},
			Assertions: []grill.Assertion{
				helper.AssertItemsCount("test_index", 0, nil),
			},
			Cleaners: []grill.Cleaner{
				helper.DeleteIndices("test_index"),
			},
		},
	}

	grill.Run(t, tests)
}
