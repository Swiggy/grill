package grill

import (
	"testing"
)

func Run(t *testing.T, testCases []TestCase) {
	for _, testCase := range testCases {
		testCase.Run(t)
	}
}

func RunParallel(t *testing.T, testCases []TestCase) {
	for _, testCase := range testCases {
		go testCase.Run(t)
	}
}
