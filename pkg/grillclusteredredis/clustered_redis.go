package grillclusterredis

import (
	"context"
	"github.com/Swiggy/grill/canned"
	"github.com/go-redis/redis/v8"
)

type ClusteredRedis struct {
	clusteredRedis *canned.ClusteredRedis
}

func (cr *ClusteredRedis) GetClusteredRedisClient() redis.UniversalClient {
	return cr.clusteredRedis.RedisClient
}

func (cr *ClusteredRedis) Start(ctx context.Context) error {
	clusterRedis, err := canned.NewClusteredRedis()
	cr.clusteredRedis = clusterRedis
	if err != nil {
		return err
	}
	return nil
}

func (cr *ClusteredRedis) Host() string {
	return cr.clusteredRedis.Host
}

func (cr *ClusteredRedis) Port() string {
	return cr.clusteredRedis.Port
}

func (cr *ClusteredRedis) Stop(ctx context.Context) error {
	return cr.clusteredRedis.Container.Terminate(ctx)
}

func (gr *ClusteredRedis) Get(key string) (string, error) {
	conn := gr.GetClusteredRedisClient()
	return conn.Get(context.Background(), key).Result()
}

func (gr *ClusteredRedis) HGetAll(key string) (map[string]string, error) {
	conn := gr.GetClusteredRedisClient()
	hgetAllResponse := conn.HGetAll(context.Background(), key)
	return hgetAllResponse.Val(), hgetAllResponse.Err()
}
