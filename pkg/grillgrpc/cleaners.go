package grillgrpc

import "github.com/Swiggy/grill"

func (gg *GRPC) ResetAllStubs() grill.Cleaner {
	return grill.CleanerFunc(func() error {
		gg.recorder.resetAll()
		return nil
	})
}
