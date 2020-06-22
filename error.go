package grill

import "fmt"

func AssertError(assertion Assertion) Assertion {
	return AssertionFunc(func() error {
		if err := assertion.Assert(); err == nil {
			return fmt.Errorf("expected error in assertion=%v, got=%v", assertion, false)
		}
		return nil
	})
}
