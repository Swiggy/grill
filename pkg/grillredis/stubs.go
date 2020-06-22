package grillredis

import "bitbucket.org/swigy/grill"

func (gr *Redis) SelectDB(db int) grill.Stub {
	return grill.StubFunc(func() error {
		_, err := gr.Client().Do("select", db)
		return err
	})
}

func (gr *Redis) Set(key, value string) grill.Stub {
	return grill.StubFunc(func() error {
		_, err := gr.Client().Do("SET", key, value)
		return err
	})
}
