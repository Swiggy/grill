package grilltile38

import (
	"bitbucket.org/swigy/grill"
)

func (gt *Tile38) SetObject(key string, id string, object string) grill.Stub {
	return grill.StubFunc(func() error {
		_, err := gt.Client().Do("SET", key, id, "OBJECT", object)
		return err
	})
}
