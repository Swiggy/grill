package grillredis

import "bitbucket.org/swigy/grill"

func (gr *Redis) FlushDB() grill.Cleaner {
	return grill.CleanerFunc(func() error {
		_, err := gr.Client().Do("flushdb")
		return err
	})
}
