package grillredis

import (
	"context"

	"bitbucket.org/swigy/grill/internal/canned"
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

func (gr *Redis) Client() redis.Conn {
	return gr.redis.Client
}

func (gr *Redis) Stop() error {
	return gr.redis.Container.Terminate(context.Background())
}
