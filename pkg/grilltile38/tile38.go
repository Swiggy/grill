package grilltile38

import (
	"context"

	"bitbucket.org/swigy/grill/internal/canned"
	"github.com/gomodule/redigo/redis"
)

type Tile38 struct {
	tile38 *canned.Tile38
}

func (gt *Tile38) Start(ctx context.Context) error {
	tile38, err := canned.NewTile38(ctx)
	if err != nil {
		return err
	}

	gt.tile38 = tile38
	return nil
}

func (gt *Tile38) Host() string {
	return gt.tile38.Host
}

func (gt *Tile38) Port() string {
	return gt.tile38.Port
}

func (gt *Tile38) Client() redis.Conn {
	return gt.tile38.Client
}

func (gt *Tile38) Stop() error {
	return gt.tile38.Container.Terminate(context.Background())
}