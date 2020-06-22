package grillgrpc

import (
	"fmt"

	"bitbucket.org/swigy/grill"
)

func (gg *GRPC) AssertCount(request Request, expectedCount int) grill.Assertion {
	return grill.AssertionFunc(func() error {
		got := gg.recorder.count(&request)
		if got != expectedCount {
			return fmt.Errorf("invalid number of requests, got=%v, want=%v", got, expectedCount)
		}
		return nil
	})
}
