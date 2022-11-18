package grillredis

import (
	"context"
	"github.com/Swiggy/grill/canned"

	"github.com/gomodule/redigo/redis"
)

type Redis struct {
	redis *canned.Redis
}

func (gr *Redis) Start(ctx context.Context) error {
	redis, err := canned.NewRedis(ctx)
	if err != nil {
		return err
	}
	gr.redis = redis
	return nil
}

func (gr *Redis) Host() string {
	return gr.redis.Host
}

func (gr *Redis) Port() string {
	return gr.redis.Port
}

func (gr *Redis) Pool() *redis.Pool {
	return gr.redis.Pool
}

func (gr *Redis) Stop(ctx context.Context) error {
	return gr.redis.Container.Terminate(ctx)
}
