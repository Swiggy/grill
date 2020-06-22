package grillkafka

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"bitbucket.org/swigy/grill"
)

func (gk *Kafka) AssertCount(topicName string, expectedCount int) grill.Assertion {
	return grill.AssertionFunc(func() error {
		group := fmt.Sprintf("%s_%d_%s", "oh_my_test_helper", rand.Intn(1000), time.Now())
		consumer, err := gk.NewConsumer(group, topicName, time.Second*5)
		if err != nil {
			return err
		}

		count := 0
		for _ = range consumer.MsgChan() {
			count++
		}

		if count != expectedCount {
			return fmt.Errorf("invalid number of messages, got=%v, want=%v", count, expectedCount)
		}

		return nil
	})
}

func (gk *Kafka) AssertMessagePresent(message Message) grill.Assertion {
	return grill.AssertionFunc(func() error {
		group := fmt.Sprintf("%s_%d_%s", "oh_my_test_helper", rand.Intn(1000), time.Now())
		consumer, err := gk.NewConsumer(group, message.Topic, time.Second*5)
		if err != nil {
			return err
		}

		for msg := range consumer.MsgChan() {
			if reflect.DeepEqual(msg, message) {
				return nil
			}
		}
		return fmt.Errorf("message not found, message=%v", message)
	})
}
