package grillconsul

import (
	"fmt"

	"bitbucket.org/swigy/grill"
)

func (gc *GrillConsul) AssertValue(key, expected string) grill.Assertion {
	return grill.AssertionFunc(func() error {
		pair, _, err := gc.consul.Client.KV().Get(key, nil)
		if err != nil {
			return err
		}
		if pair == nil {
			return fmt.Errorf("no value found for key=%v", key)
		}
		if string(pair.Value) != expected {
			return fmt.Errorf("invalid value for key=%v, got=%v, want=%v", key, string(pair.Value), expected)
		}
		return nil
	})
}
