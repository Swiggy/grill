package grills3

import (
	"context"
	"testing"

	"github.com/Swiggy/grill"
)

func Test_GrillS3(t *testing.T) {
	helper := &S3{}
	if err := helper.Start(context.TODO()); err != nil {
		t.Errorf("error starting s3 grill, error=%v", err)
		return
	}

	tests := []grill.TestCase{
		{
			Name: "Test_UploadExists",
			Stubs: []grill.Stub{
				helper.CreateBucket("test-bucket"),
				helper.UploadFile("test-bucket", "testFile.txt", "test_data/testFile.txt"),
			},
			Action: func() interface{} {
				return nil
			},
			Assertions: []grill.Assertion{
				helper.AssertFileExists("test-bucket", "testFile.txt"),
			},
			Cleaners: []grill.Cleaner{
				helper.DeleteAllFiles("test-bucket"),
				helper.DeleteBucket("test-bucket"),
			},
		},
	}

	grill.Run(t, tests)
}
