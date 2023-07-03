package grillclusterredis

import (
	"context"
	"fmt"
	"github.com/Swiggy/grill"
)

func (gr *ClusteredRedis) AssertValue(key, expected string) grill.Assertion {
	return grill.AssertionFunc(func() error {
		conn := gr.GetClusteredRedisClient()

		output := conn.Get(context.Background(), key)
		if output.Err() != nil {
			return output.Err()
		}
		if output == nil {
			return fmt.Errorf("no value found for key=%v", key)
		}
		got := output.Val()
		if got != expected {
			return fmt.Errorf("invalid value for key=%v, got=%v, want=%v", key, got, expected)
		}
		return nil
	})
}
