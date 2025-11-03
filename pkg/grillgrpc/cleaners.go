package grillgrpc

import "github.com/singh-jatin28/grill"

func (gg *GRPC) ResetAllStubs() grill.Cleaner {
	return grill.CleanerFunc(func() error {
		gg.recorder.resetAll()
		return nil
	})
}
