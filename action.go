package grill

import (
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
)

var Any = struct{}{}

type OutputAssertion interface {
	Assertion
	SetOutput(output interface{})
}

func ActionOutput(output ...interface{}) interface{} {
	return output
}

func AssertOutput(args ...interface{}) Assertion {
	return &assertOutput{expected: args}
}

type assertOutput struct {
	output   []interface{}
	expected []interface{}
}

func (assert *assertOutput) SetOutput(output interface{}) {
	if slice, ok := output.([]interface{}); ok {
		assert.output = slice
		return
	}
	assert.output = []interface{}{output}
}

func (assert *assertOutput) Assert() error {
	if len(assert.output) != len(assert.expected) {
		return fmt.Errorf("invalid number of arguments in expected, OutputLength=%v, ExpectedLength=%v", len(assert.output), len(assert.expected))
	}

	for i := 0; i < len(assert.output); i++ {
		if assert.expected[i] != Any {
			if out, ok := assert.output[i].(proto.Message); ok {
				if assert.expected[i] == nil && reflect.ValueOf(assert.output[i]).IsNil() {
					// added to check for nil proto message, as proto.Equal does not work for nil proto messages
					continue
				}
				exp, _ := assert.expected[i].(proto.Message)
				if !proto.Equal(out, exp) {
					return fmt.Errorf("invalid proto message value at index=%v, got=%v, want=%v", i, assert.output[i], assert.expected[i])
				}
			} else {
				if !reflect.DeepEqual(assert.output[i], assert.expected[i]) {
					return fmt.Errorf("invalid value at index=%v, got=%v, want=%v", i, assert.output[i], assert.expected[i])
				}
			}
		}
	}
	return nil
}
