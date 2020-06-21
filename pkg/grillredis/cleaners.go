package grillredis

import "bitbucket.org/swigy/grill"

func (gr *GrillRedis) FlushDB() grill.Cleaner {
	return grill.CleanerFunc(func() error {
		_, err := gr.Client().Do("flushdb")
		return err
	})
}
