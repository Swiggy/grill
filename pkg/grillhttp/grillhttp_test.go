package grillhttp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"bitbucket.org/swigy/grill"
)

var (
	testStub = Stub{
		Request: Request{
			Method:  "GET",
			UrlPath: "/test",
		},
		Response: Response{
			Status: 200,
			Body:   "PASS",
		},
	}
)

func Test_GrillHTTP(t *testing.T) {
	helper, err := Start()
	if err != nil {
		t.Errorf("Error Creating Stub, error = %v", err)
		return
	}

	tests := []grill.TestCase{
		{
			Name:  "TestMockHTTP_StubNotPresent",
			Stubs: []grill.Stub{},
			Action: func() interface{} {
				res, err := http.Get(fmt.Sprintf("http://localhost:%s/test", helper.Port()))
				if res == nil || res.Body == nil {
					return grill.MultiOutput(nil, err)
				}
				defer res.Body.Close()
				got, _ := ioutil.ReadAll(res.Body)

				return grill.MultiOutput(string(got), res.StatusCode, err)
			},
			Assertions: []grill.Assertion{
				grill.AssertOutput(grill.Any, 404, nil),
				helper.AssertCount(&testStub.Request, 1),
			},
			Cleaners: []grill.Cleaner{
				helper.ResetAllStubs(),
			},
		},
		{
			Name: "TestMockHTTP_StubPresent",
			Stubs: []grill.Stub{
				helper.Stub(&testStub),
			},
			Action: func() interface{} {
				res, err := http.Get(fmt.Sprintf("http://localhost:%s/test", helper.Port()))
				if res == nil || res.Body == nil {
					return grill.MultiOutput(nil, err)
				}
				defer res.Body.Close()
				got, _ := ioutil.ReadAll(res.Body)

				return grill.MultiOutput(string(got), res.StatusCode, err)
			},
			Assertions: []grill.Assertion{
				grill.AssertOutput("PASS", 200, nil),
				helper.AssertCount(&testStub.Request, 1),
			},
			Cleaners: []grill.Cleaner{
				helper.ResetAllStubs(),
			},
		},
	}

	grill.Run(t, tests)
}
