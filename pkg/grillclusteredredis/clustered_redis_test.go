package grillclusterredis

import (
	"context"
	"fmt"
	"github.com/Swiggy/grill"
	"testing"
)

func Test_GrillClusteredRedis(t *testing.T) {
	helper := &ClusteredRedis{}
	if err := helper.Start(context.TODO()); err != nil {
		t.Errorf("error starting redis grill, error=%v", err)
		return
	}
	const TotalNumberOfTestCasesToRun = 100

	tests := []grill.TestCase{
		{
			Name:  "Test_GetSet-AllShards",
			Stubs: []grill.Stub{},
			Action: func() interface{} {
				for testCaseCounter := 0; testCaseCounter <= TotalNumberOfTestCasesToRun; testCaseCounter++ {
					err := helper.Set(fmt.Sprintf("key-%d", testCaseCounter), fmt.Sprintf("val-%d", testCaseCounter), 0).Stub()
					if err != nil {
						return grill.ActionOutput(err)
					}
				}
				return grill.ActionOutput(nil)
			},
			Assertions: []grill.Assertion{
				grill.AssertOutput(nil),
				grill.AssertionFunc(func() error {
					for testCaseCounter := 0; testCaseCounter <= TotalNumberOfTestCasesToRun; testCaseCounter++ {
						err := helper.AssertValue(fmt.Sprintf("key-%d", testCaseCounter), fmt.Sprintf("val-%d", testCaseCounter)).Assert()
						if err != nil {
							return err
						}
					}
					return nil
				}),
			},
			Cleaners: []grill.Cleaner{
				helper.FlushDB(),
			},
		},
	}
	grill.Run(t, tests)
	if err := helper.Stop(context.TODO()); err != nil {
		t.Errorf("error starting redis grill, error=%v", err)
		return
	}

}
