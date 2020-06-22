package grillgrpc

import "bitbucket.org/swigy/grill"

func (gg *GrillGRPC) ResetAllStubs() grill.Cleaner {
	return grill.CleanerFunc(func() error {
		gg.recorder.resetAll()
		return nil
	})
}
