package grillgrpc

import (
	"bitbucket.org/swigy/grill"
)

func (gg *GRPC) Stub(stub *Stub) grill.Stub {
	return grill.StubFunc(func() error {
		return gg.recorder.add(stub)
	})
}
