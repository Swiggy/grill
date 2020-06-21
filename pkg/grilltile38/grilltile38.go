package grilltile38

import (
	"context"

	"bitbucket.org/swigy/grill/internal/canned"
	"github.com/gomodule/redigo/redis"
)

type GrillTile38 struct {
	tile38 *canned.Tile38
}

func Start() (*GrillTile38, error) {
	tile38, err := canned.NewTile38(context.TODO())
	if err != nil {
		return nil, err
	}

	return &GrillTile38{
		tile38: tile38,
	}, nil
}

func (gt *GrillTile38) Host() string {
	return gt.tile38.Host
}

func (gt *GrillTile38) Port() string {
	return gt.tile38.Port
}

func (gt *GrillTile38) Client() redis.Conn {
	return gt.tile38.Client
}

func (gt *GrillTile38) Stop() error {
	return gt.tile38.Container.Terminate(context.Background())
}
