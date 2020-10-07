package grill

import (
	"sync"
	"testing"
)

func Run(t *testing.T, testCases []TestCase) {
	for _, testCase := range testCases {
		testCase.Run(t)
	}
}

func RunParallel(t *testing.T, testCases []TestCase) {
	wg := &sync.WaitGroup{}
	wg.Add(len(testCases))

	for _, testCase := range testCases {
		go func(tt TestCase, wg *sync.WaitGroup) {
			defer wg.Done()

			testCase.Run(t)
		}(testCase, wg)
	}

	wg.Wait()
}
