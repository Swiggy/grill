package grilltile38

import (
	"context"
	"testing"

	"github.com/Swiggy/grill"
)

var (
	data  = `{"type":"Polygon","coordinates":[[[77.703033,12.941534],[77.703297,12.942202],[77.70362,12.942083],[77.703373,12.941421],[77.703037,12.941532],[77.703033,12.941534]]]}`
	data2 = `{"type":"Polygon","coordinates":[[[77.7030331,12.9415341],[77.703297,12.942202],[77.70362,12.942083],[77.703373,12.941421],[77.703037,12.941532],[77.7030331,12.9415341]]]}`
)

func Test_GrillTile38(t *testing.T) {
	helper := &Tile38{}
	if err := helper.Start(context.TODO()); err != nil {
		t.Errorf("error starting tile38 grill, error=%v", err)
		return
	}
	tests := []grill.TestCase{
		{
			Name: "Test_SetGet",
			Stubs: []grill.Stub{
				helper.SetObject("one", "1", data),
				helper.SetObject("one", "2", data2),
			},
			Action: func() interface{} {
				return nil
			},
			Assertions: []grill.Assertion{
				helper.AssertObject("one", "1", data),
				helper.AssertObject("one", "1", data),
				helper.AssertObject("one", "2", data2),
				helper.AssertObject("one", "1", data),
			},
			Cleaners: []grill.Cleaner{
				helper.FlushDB(),
			},
		},
		{
			Name: "Test_Exist",
			Stubs: []grill.Stub{
				helper.SetObject("one", "1", data),
			},
			Action: func() interface{} {
				return nil
			},
			Assertions: []grill.Assertion{
				helper.AssertExist("one", "1"),
			},
			Cleaners: []grill.Cleaner{
				helper.FlushDB(),
			},
		},
		{
			Name: "Test_NotExist",
			Stubs: []grill.Stub{
				helper.SetObject("one", "1", data),
			},
			Action: func() interface{} {
				return nil
			},
			Assertions: []grill.Assertion{
				helper.AssertExist("one", "1"),
				grill.Not(helper.AssertExist("onee", "2")),
			},
			Cleaners: []grill.Cleaner{
				helper.FlushDB(),
			},
		},
	}

	grill.Run(t, tests)
}
