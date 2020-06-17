package grill

import (
	"sync"
	"testing"
)

func Run(t *testing.T, testCases []TestCase) {
	for _, testCase := range testCases {
		run(t, testCase)
	}
}

func run(t *testing.T, testCase TestCase) {
	t.Run(testCase.Name, func(t *testing.T) {
		defer func() {
			for _, cleaner := range testCase.Cleaners {
				if err := cleaner.Clean(); err != nil {
					t.Errorf("error cleaning stub, error=%v", err)
				}
			}
		}()

		for _, stub := range testCase.Stubs {
			if err := stub.Stub(); err != nil {
				t.Errorf("error stub, error=%v", err)
				return
			}
		}

		output := testCase.Action()

		wg := sync.WaitGroup{}
		wg.Add(len(testCase.Assertions))

		for _, a := range testCase.Assertions {
			go func(wg *sync.WaitGroup, assertion Assertion) {
				defer wg.Done()
				if assertion, ok := assertion.(OutputAssertion); ok {
					assertion.SetOutput(output)
					if err := assertion.Assert(); err != nil {
						t.Errorf("assertion failed, error=%v", err)
					}
					return
				}
				if err := assertion.Assert(); err != nil {
					t.Errorf("assertion failed, error=%v", err)
				}
			}(&wg, a)
		}

		wg.Wait()
	})
}
