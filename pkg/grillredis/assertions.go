package grillredis

import (
	"fmt"

	"bitbucket.org/swigy/grill"
)

func (gr *Redis) AssertValue(key, expected string) grill.Assertion {
	return grill.AssertionFunc(func() error {
		conn := gr.Pool().Get()
		defer conn.Close()

		output, err := conn.Do("GET", key)
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
