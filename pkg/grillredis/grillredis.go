package grillredis

import (
	"context"

	"bitbucket.org/swigy/grill/internal/canned"
	"github.com/gomodule/redigo/redis"
)

type GrillRedis struct {
	redis *canned.Redis
}

func Start() (*GrillRedis, error) {
	redis, err := canned.NewRedis(context.TODO())
	if err != nil {
		return nil, err
	}

	return &GrillRedis{
		redis: redis,
	}, nil
}

func (gr *GrillRedis) Host() string {
	return gr.redis.Host
}

func (gr *GrillRedis) Port() string {
	return gr.redis.Port
}

func (gr *GrillRedis) Client() redis.Conn {
	return gr.redis.Client
}

func (gr *GrillRedis) Stop() error {
	return gr.redis.Container.Terminate(context.Background())
}
