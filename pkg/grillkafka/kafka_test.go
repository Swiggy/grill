package grillkafka

import (
	"fmt"
	"testing"
	"time"

	"bitbucket.org/swigy/grill"
)

var (
	testTopic      = "test_topic"
	testErrorTopic = "test_topic_error"

	testMessage = Message{
		Key:     "1",
		Headers: map[string]string{},
		Value:   "abc",
	}
)

func Test_GrillKafka(t *testing.T) {
	helper, err := Start()
	if err != nil {
		t.Errorf("error starting kafka grill, error=%v", err)
	}

	test := grill.TestCase{
		Name: "Test_GrillKafka_ProduceConsume",
		Stubs: []grill.Stub{
			helper.CreateTopics("test_topic", "test_topic_error"),
		},
		Action: func() interface{} {
			err := helper.Produce(testTopic, testMessage)
			if err != nil {
				return false
			}
			go func() {
				consumer, err := helper.NewConsumer("test_consumer", testTopic, time.Second*5)
				if err != nil {
					return
				}
				for msg := range consumer.MsgChan() {
					if err := helper.Produce(testErrorTopic, msg); err != nil {
						fmt.Printf("error producing to error topic, error=%v\n", err)
					}
				}
			}()
			return true
		},
		Assertions: []grill.Assertion{
			grill.AssertOutput(true),
			grill.Try(time.Second*30, 3, helper.AssertCount("test_topic", 1)),
			grill.Try(time.Second*30, 3, helper.AssertMessageCount("test_topic_error", testMessage, 1)),
		},
		Cleaners: []grill.Cleaner{
			helper.DeleteTopics("test_topic", "test_topic_error"),
		},
	}

	grill.Run(t, []grill.TestCase{test, test})
}
