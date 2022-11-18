package grilltile38

import "github.com/Swiggy/grill"

func (gt *Tile38) FlushDB() grill.Cleaner {
	return grill.CleanerFunc(func() error {
		conn := gt.Pool().Get()
		defer conn.Close()

		_, err := conn.Do("flushdb")
		return err
	})
}
