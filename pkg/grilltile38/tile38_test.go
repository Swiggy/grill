package grilltile38

import (
	"testing"

	"bitbucket.org/swigy/grill"
)

var (
	data = `{"type":"Polygon","coordinates":[[[77.703033,12.941534],[77.703297,12.942202],[77.70362,12.942083],[77.703373,12.941421],[77.703037,12.941532],[77.703033,12.941534]]]}`
)

func Test_GrillTile38(t *testing.T) {
	helper, err := Start()
	if err != nil {
		t.Errorf("error starting tile38 grill, error=%v", err)
		return
	}

	tests := []grill.TestCase{
		{
			Name: "Test_GetSet",
			Stubs: []grill.Stub{
				helper.SetObject("one", "1", data),
			},
			Action: func() interface{} {
				return nil
			},
			Assertions: []grill.Assertion{
				helper.AssertObject("one", "1", data),
			},
			Cleaners: []grill.Cleaner{
				helper.FlushDB(),
			},
		},
	}

	grill.Run(t, tests)
}
