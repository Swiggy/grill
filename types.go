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
