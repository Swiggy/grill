package grilltile38

import (
	"github.com/singh-jatin28/grill"
)

func (gt *Tile38) SetObject(key string, id string, object string) grill.Stub {
	return grill.StubFunc(func() error {
		conn := gt.Pool().Get()
		defer conn.Close()

		_, err := conn.Do("SET", key, id, "OBJECT", object)
		return err
	})
}
