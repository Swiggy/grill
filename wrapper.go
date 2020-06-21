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

func StubFunc(fn func() error) Stub {
	return wrapper{fn: fn}
}

func AssertionFunc(fn func() error) Assertion {
	return wrapper{fn: fn}
}

func CleanerFunc(fn func() error) Cleaner {
	return wrapper{fn: fn}
}
