package grill

import (
	"fmt"
	"time"
)

func Try(deadline time.Duration, minSuccess int, assertion Assertion) Assertion {
	return &tryAssertion{
		assertion:  assertion,
		deadline:   deadline,
		minSuccess: minSuccess,
	}
}

type tryAssertion struct {
	assertion  Assertion
	deadline   time.Duration
	minSuccess int
}

func (assert *tryAssertion) Assert() error {
	checkC := time.Tick(assert.deadline / time.Duration(assert.minSuccess*2+2))
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
			if successCount >= assert.minSuccess {
				return nil
			}
		case <-quitC:
			return fmt.Errorf("couldn't complete in given deadline, max consecutive success=%d", successCount)
		}
	}
	return nil
}
