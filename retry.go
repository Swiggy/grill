package grill

import (
	"fmt"
	"time"
)

var Retry = func(assertion Assertion, deadline time.Duration, success int) Assertion {
	return &retryAssertion{
		assertion: assertion,
		deadline:  deadline,
		success:   success,
	}
}

type retryAssertion struct {
	assertion Assertion
	deadline  time.Duration
	success   int
}

func (assert *retryAssertion) Assert() error {
	checkC := time.Tick(assert.deadline / 11)
	quitC := time.Tick(assert.deadline)
	var successCount = 0
	for {
		select {
		case <-checkC:
			if err := assert.assertion.Assert(); err != nil {
				successCount = 0
				continue
			}
			successCount += 1
			if successCount >= assert.success {
				return nil
			}
		case <-quitC:
			return fmt.Errorf("couldn't complete in given deadline, max consecutive success=%d", successCount)
		}
	}
	return nil
}
