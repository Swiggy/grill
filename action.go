package grill

import (
	"fmt"
	"reflect"
)

var Any = struct{}{}

type OutputAssertion interface {
	Assertion
	SetOutput(output interface{})
}

var MultiOutput = func(output ...interface{}) interface{} {
	return output
}

var AssertOutput = func(args ...interface{}) Assertion {
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
			if !reflect.DeepEqual(assert.output[i], assert.expected[i]) {
				return fmt.Errorf("invalid value at index=%v, got=%v, want=%v", i, assert.output[i], assert.expected[i])
			}
		}
	}
	return nil
}
