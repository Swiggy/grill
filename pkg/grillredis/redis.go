package grillredis

import (
	"context"

	"bitbucket.org/swigy/grill/internal/canned"
	"github.com/gomodule/redigo/redis"
)

type Redis struct {
	redis *canned.Redis
}

func Start() (*Redis, error) {
	redis, err := canned.NewRedis(context.TODO())
	if err != nil {
		return nil, err
	}

	return &Redis{
		redis: redis,
	}, nil
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
