package grillconsul

import (
	"context"
	"testing"

	"github.com/Swiggy/grill"
)

func Test_GrillRedis(t *testing.T) {
	helper := &Consul{}
	if err := helper.Start(context.TODO()); err != nil {
		t.Errorf("error starting consul grill, error=%v", err)
		return
	}

	tests := []grill.TestCase{
		{
			Name: "Test_Seed",
			Stubs: []grill.Stub{
				helper.SeedFromCSVFile("test_data/seed.csv"),
			},
			Action: func() interface{} {
				return nil
			},
			Assertions: []grill.Assertion{
				helper.AssertValue("test/one", "1"),
				helper.AssertValue("test/two", "2"),
				helper.AssertValue("test/three", "3"),
			},
			Cleaners: []grill.Cleaner{
				helper.DeleteAllKeys(),
			},
		},
		{
			Name: "Test_SetGet",
			Stubs: []grill.Stub{
				helper.Set("test/four", "4"),
			},
			Action: func() interface{} {
				return nil
			},
			Assertions: []grill.Assertion{
				helper.AssertValue("test/four", "4"),
			},
			Cleaners: []grill.Cleaner{
				helper.DeleteAllKeys(),
			},
		},
	}

	grill.Run(t, tests)
}
