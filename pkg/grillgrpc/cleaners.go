package grillgrpc

import "github.com/lovlin-thakkar/swiggy-grill"

func (gg *GRPC) ResetAllStubs() grill.Cleaner {
	return grill.CleanerFunc(func() error {
		gg.recorder.resetAll()
		return nil
	})
}
