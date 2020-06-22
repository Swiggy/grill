package grillgrpc

import "bitbucket.org/swigy/grill"

func (gg *GRPC) ResetAllStubs() grill.Cleaner {
	return grill.CleanerFunc(func() error {
		gg.recorder.resetAll()
		return nil
	})
}
