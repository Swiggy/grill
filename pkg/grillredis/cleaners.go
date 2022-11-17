package grillredis

import "github.com/lovlin-thakkar/swiggy-grill"

func (gr *Redis) FlushDB() grill.Cleaner {
	return grill.CleanerFunc(func() error {
		conn := gr.Pool().Get()
		defer conn.Close()

		_, err := conn.Do("flushdb")
		return err
	})
}
