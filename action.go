package grill

var Any = struct{}{}

type actionOutput struct {
	output []interface{}
}

var ActionOutput = func(out ...interface{}) interface{} {
	return actionOutput{output: out}
}

var AssertOutput = func(expected ...interface{}) Assertion {
	return &assertOutput{Expected: expected}
}

type assertOutput struct {
	output   interface{}
	Expected []interface{}
}

func (assert *assertOutput) Assert() error {
	return nil
}
