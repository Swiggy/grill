package grill

import (
	"sync"
	"testing"
)

type TestCase struct {
	Name       string
	Stubs      []Stub
	Action     func() interface{}
	Assertions []Assertion
	Cleaners   []Cleaner
}

func (testCase *TestCase) Run(t *testing.T) {
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
