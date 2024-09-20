package grill

import (
	"reflect"
	"testing"

	"github.com/golang/protobuf/ptypes/wrappers"
)

func TestMultiOutput(t *testing.T) {
	tests := []struct {
		name  string
		input []interface{}
		want  interface{}
	}{
		{"NoElement", []interface{}{}, []interface{}{}},
		{"OneElement", []interface{}{1}, []interface{}{1}},
		{"MultiElement", []interface{}{1, 2}, []interface{}{1, 2}},
		{"CompositeElements", []interface{}{1, "abc", nil}, []interface{}{1, "abc", nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ActionOutput(tt.input...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ActionOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssertOutput(t *testing.T) {
	type custom struct {
		val int
	}
	tests := []struct {
		name     string
		output   interface{}
		expected []interface{}
		wantErr  bool
	}{
		{"SingleInt-Failure", 1, []interface{}{2}, true},
		{"SingleInt-Success", 1, []interface{}{1}, false},
		{"MultiInt-Failure", []interface{}{1, 2, 3}, []interface{}{1, 1, 3}, true},
		{"MultiIntInt-InvalidLength-Failure", []interface{}{1, 2, 3, 4}, []interface{}{1, 1, 3}, true},
		{"MultiInt-AnySuccess", []interface{}{1, 2, 3}, []interface{}{Any, Any, Any}, false},
		{"MultiInt-Success", []interface{}{1, 2, 3}, []interface{}{1, 2, 3}, false},
		{"String-Success", []interface{}{"1", "2", "3"}, []interface{}{"1", "2", "3"}, false},
		{"Custom-Success", &custom{1}, []interface{}{&custom{1}}, false},
		{"Composite-Failure", []interface{}{1, "1", &custom{1}}, []interface{}{1, "1", &custom{2}}, true},
		{"Composite-Success", []interface{}{1, "1", &custom{1}}, []interface{}{1, "1", &custom{1}}, false},
		{"Composite-Any", []interface{}{1, "1", &custom{1}}, []interface{}{1, "1", Any}, false},
		{"NilProtoMessage-Success", []interface{}{getNilProtoMessage(), "1"}, []interface{}{nil, "1"}, false},
		{"NilError-Success", []interface{}{2, getNilError()}, []interface{}{2, nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertion := AssertOutput(tt.expected...)
			assertion.(OutputAssertion).SetOutput(tt.output)
			if gotErr := assertion.Assert(); (gotErr != nil) != tt.wantErr {
				t.Errorf("AssertOutput() gotErr = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}

func getNilProtoMessage() *wrappers.StringValue {
	return nil
}

func getNilError() error {
	return nil
}
