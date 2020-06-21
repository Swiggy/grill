package grilltile38

import (
	"encoding/json"
	"fmt"
	"reflect"

	"bitbucket.org/swigy/grill"
)

func (gt *GrillTile38) AssertObject(key, id string, expected string) grill.Assertion {
	return grill.AssertionFunc(func() error {
		output, err := gt.Client().Do("GET", key, id)
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
