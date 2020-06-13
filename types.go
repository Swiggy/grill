package grill

type Stub interface {
	Stub() error
}

type Assertion interface {
	Assert() error
}

type Cleaner interface {
	Clean() error
}

type TestCase struct {
	Name       string
	Stubs      []Stub
	Action     func() interface{}
	Assertions []Assertion
	Cleaners   []Cleaner
}
