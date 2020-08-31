package grillredis

import "bitbucket.org/swigy/grill"

func (gr *Redis) SelectDB(db int) grill.Stub {
	return grill.StubFunc(func() error {
		conn := gr.Pool().Get()
		defer conn.Close()

		_, err := conn.Do("select", db)
		return err
	})
}

func (gr *Redis) Set(key, value string) grill.Stub {
	return grill.StubFunc(func() error {
		conn := gr.Pool().Get()
		defer conn.Close()

		_, err := conn.Do("SET", key, value)
		return err
	})
}
