package grill

import (
	"fmt"
	"time"
)

type TryOption func(*tryAssertion)

// WithCheckFrequency sets a custom frequency for assertion checks.
// If not set, frequency is derived from deadline and minSuccess.
func WithCheckFrequency(frequency time.Duration) TryOption {
	return func(t *tryAssertion) {
		t.checkFrequency = frequency
	}
}

func Try(deadline time.Duration, minSuccess int, assertion Assertion, options ...TryOption) Assertion {
	t := &tryAssertion{
		assertion:      assertion,
		deadline:       deadline,
		minSuccess:     minSuccess,
		checkFrequency: deadline / time.Duration(minSuccess*3+3),
	}
	for _, opt := range options {
		opt(t)
	}
	return t
}

type tryAssertion struct {
	assertion      Assertion
	deadline       time.Duration
	minSuccess     int
	checkFrequency time.Duration
}

func (assert *tryAssertion) Assert() error {
	checkC := time.Tick(assert.checkFrequency)
	quitC := time.Tick(assert.deadline)
	var successCount = 0
	var errors []string
	for {
		select {
		case <-checkC:
			if err := assert.assertion.Assert(); err != nil {
				errors = append(errors, err.Error())
				successCount = 0
				continue
			}
			successCount += 1
			if successCount >= assert.minSuccess {
				return nil
			}
		case <-quitC:
			return fmt.Errorf("couldn't complete in given deadline, max consecutive success=%d, errors=%v", successCount, uniq(errors))
		}
	}
}

func uniq(errors []string) []string {
	temp := map[string]struct{}{}
	var result []string
	for _, e := range errors {
		if _, ok := temp[e]; !ok {
			result = append(result, e)
			temp[e] = struct{}{}
		}
	}
	return result
}
