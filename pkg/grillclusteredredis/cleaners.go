package grillclusterredis

import (
	"context"
	"github.com/Swiggy/grill"
)

func (gr *ClusteredRedis) FlushDB() grill.Cleaner {
	return grill.CleanerFunc(func() error {
		conn := gr.GetClusteredRedisClient()

		output := conn.FlushDB(context.Background())
		return output.Err()
	})
}
