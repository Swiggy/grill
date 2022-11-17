package grilltile38

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gomodule/redigo/redis"
	"github.com/lovlin-thakkar/swiggy-grill"
)

func (gt *Tile38) AssertObject(key, id string, expected string) grill.Assertion {
	return grill.AssertionFunc(func() error {
		conn := gt.Pool().Get()
		defer conn.Close()

		output, err := conn.Do("GET", key, id)
		if err != nil {
			return err
		}
		if output == nil {
			return fmt.Errorf("no value found for key=%v", key)
		}
		obj := map[string]interface{}{}
		err = json.Unmarshal(output.([]byte), &obj)
		if e, ok := obj["err"]; ok {
			return fmt.Errorf("%s", e)
		}

		want := map[string]interface{}{}
		if err := json.Unmarshal([]byte(expected), &want); err != nil {
			return fmt.Errorf("error unmarshalling expected object, error=%v", err)
		}
		if !reflect.DeepEqual(obj["object"], want) {
			return fmt.Errorf("invalid value for key=%v, id=%v, got=%v, want=%v", key, id, obj["object"], want)
		}
		return nil
	})
}

func (gt *Tile38) AssertExist(key, id string) grill.Assertion {
	return grill.AssertionFunc(func() error {
		conn := gt.Pool().Get()
		defer conn.Close()

		output, err := redis.String(conn.Do("GET", key, id))
		if err != nil {
			return err
		}
		if output == "" {
			return fmt.Errorf("no value found for key=%v", key)
		}
		obj := map[string]interface{}{}
		if err := json.Unmarshal([]byte(output), &obj); err != nil {
			return fmt.Errorf("error unmarshalling expected object, error=%v", err)
		}
		if e, ok := obj["err"]; ok {
			return fmt.Errorf("nothing found for key=%v and id=%v, err=%v", key, id, e)
		}
		return nil
	})
}
