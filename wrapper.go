package grill

type wrapper struct {
	fn func() error
}

func (w wrapper) Stub() error {
	return w.fn()
}

func (w wrapper) Assert() error {
	return w.fn()
}

func (w wrapper) Clean() error {
	return w.fn()
}

func WrapStub(fn func() error) Stub {
	return wrapper{fn: fn}
}

func WrapAssertion(fn func() error) Assertion {
	return wrapper{fn: fn}
}

func WrapCleaner(fn func() error) Cleaner {
	return wrapper{fn: fn}
}
