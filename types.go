package grill

type Stub2 interface {
	Stub2() error
}

type Stub interface {
	Stub() error
}

type StubFunc func() error

func (fn StubFunc) Stub() error {
	return fn()
}

type Assertion interface {
	Assert() error
}

type AssertionFunc func() error

func (fn AssertionFunc) Assert() error {
	return fn()
}

type Cleaner interface {
	Clean() error
}

type CleanerFunc func() error

func (fn CleanerFunc) Clean() error {
	return fn()
}
