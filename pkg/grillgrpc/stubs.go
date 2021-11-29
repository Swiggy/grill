package grillgrpc

import (
	"github.com/swiggy-private/grill"
)

func (gg *GRPC) Stub(request Request, response Response) grill.Stub {
	return grill.StubFunc(func() error {
		return gg.recorder.add(&Stub{Request: request, Response: response})
	})
}
