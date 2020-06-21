package grillredis

import (
	"fmt"

	"bitbucket.org/swigy/grill"
)

func (gr *GrillRedis) AssertItem(key, expected string) grill.Assertion {
	return grill.AssertionFunc(func() error {
		output, err := gr.Client().Do("GET", key)
		if err != nil {
			return err
		}
		if output == nil {
			return fmt.Errorf("no value found for key=%v", key)
		}
		got := string(output.([]byte))
		if got != expected {
			return fmt.Errorf("invalid value for key=%v, got=%v, want=%v", key, got, expected)
		}
		return nil
	})
}
