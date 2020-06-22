package grills3

import (
	"testing"

	"bitbucket.org/swigy/grill"
)

func Test_GrillS3(t *testing.T) {
	helper, err := Start()
	if err != nil {
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
