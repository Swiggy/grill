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
	checkC := time.Tick(assert.deadline / time.Duration(assert.minSuccess*3+3))
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
