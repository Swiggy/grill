package grillclusterredis

import (
	"context"
	"github.com/Swiggy/grill"
	"time"
)

func (gr *ClusteredRedis) Set(key, value string, ttl time.Duration) grill.Stub {
	return grill.StubFunc(func() error {
		conn := gr.GetClusteredRedisClient()

		setResult := conn.Set(context.Background(), key, value, ttl)
		return setResult.Err()
	})
}
