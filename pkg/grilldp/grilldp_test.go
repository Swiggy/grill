package grilldp

import (
	"fmt"
	"testing"
	"time"

	dataplatform "bitbucket.org/swigy/dp-client-go"
	"bitbucket.org/swigy/grill"
)

const (
	appName   = "oh-my-test-helper"
	eventName = "test_event"
)

var (
	goodEvent = testEvent{
		Id:          "1",
		EventType:   "A",
		Version:     1,
		Score:       5.6,
		Enabled:     true,
		Tags:        []string{"one", "two"},
		TagsNonNull: []string{"three", "four"},
		Source:      testSource{Name: "tester"},
	}
)

type testSource struct {
	Name string `json:"name"`
}

type testEvent struct {
	Id          string     `json:"id"`
	EventType   string     `json:"type,omitempty"`
	Version     int        `json:"version"`
	Score       float64    `json:"score"`
	Enabled     bool       `json:"enabled"`
	Tags        []string   `json:"tags,omitempty"`
	TagsNonNull []string   `json:"tagsNonNull"`
	Source      testSource `json:"source"`
}

func newDpClient(host, port string) dataplatform.DPClient {
	return dataplatform.CreateDPClient(&dataplatform.Config{
		AppVersion:       "1.0.0",
		AppName:          appName,
		HTTPUrl:          fmt.Sprintf("http://%s:%s", host, port),
		Timeout:          time.Second,
		RetryCount:       3,
		MaxQueueSize:     100,
		AsyncBatchSize:   1,
		AsyncMaxWait:     time.Second,
		MessagePerSecond: 1,
		RetryTimeDiff:    time.Second * 10,
		LogRequests:      true,
	})
}

func Test_GrillDP(t *testing.T) {
	helper, err := Start()
	if err != nil {
		t.Errorf("error starting dp grill, error=%v", err)
		return
	}

	dpClient := newDpClient(helper.Host(), helper.Port())

	tests := []grill.TestCase{
		{
			Name: "Test_RegisteredApps",
			Action: func() interface{} {
				dpClient.RegisterEvent(eventName, "1.0.0", testEvent{})
				return nil
			},
			Assertions: []grill.Assertion{
				grill.Try(time.Second*5, 2, helper.AssertRegisteredApps(appName)),
			},
			Cleaners: []grill.Cleaner{},
		},
		{
			Name: "Test_GoodEvent",
			Action: func() interface{} {
				dpClient.SendEvent(goodEvent)
				return nil
			},
			Assertions: []grill.Assertion{
				grill.Try(time.Second*10, 2, helper.AssertCount(appName, eventName, "1.0.0", 1)),
				grill.Try(time.Second*10, 2, helper.AssertSchemaValidation(appName, eventName, "1.0.0", "test_data/schema.json")),
			},
			Cleaners: []grill.Cleaner{
				helper.FlushAllEvents(),
			},
		},
		{
			Name: "Test_NullTagsShouldFail",
			Action: func() interface{} {
				ev := goodEvent
				ev.Tags = nil
				ev.TagsNonNull = nil
				dpClient.SendEvent(ev)
				return nil
			},
			Assertions: []grill.Assertion{
				grill.Try(time.Second*10, 2, helper.AssertCount(appName, eventName, "1.0.0", 1)),
				grill.AssertError(grill.Try(time.Second*10, 2, helper.AssertSchemaValidation(appName, eventName, "1.0.0", "test_data/schema.json"))),
			},
			Cleaners: []grill.Cleaner{
				helper.FlushAllEvents(),
			},
		},
		{
			Name: "Test_InvalidEnumValueShouldFail",
			Action: func() interface{} {
				ev := goodEvent
				ev.EventType = "X"
				dpClient.SendEvent(ev)
				return nil
			},
			Assertions: []grill.Assertion{
				grill.Try(time.Second*10, 2, helper.AssertCount(appName, eventName, "1.0.0", 1)),
				grill.AssertError(grill.Try(time.Second*10, 2, helper.AssertSchemaValidation(appName, eventName, "1.0.0", "test_data/schema.json"))),
			},
			Cleaners: []grill.Cleaner{
				helper.FlushAllEvents(),
			},
		},
		{
			Name: "Test_NullTypeShouldPass",
			Action: func() interface{} {
				ev := goodEvent
				ev.EventType = ""
				dpClient.SendEvent(ev)
				return nil
			},
			Assertions: []grill.Assertion{
				grill.Try(time.Second*10, 2, helper.AssertCount(appName, eventName, "1.0.0", 1)),
				grill.Try(time.Second*10, 2, helper.AssertSchemaValidation(appName, eventName, "1.0.0", "test_data/schema.json")),
			},
			Cleaners: []grill.Cleaner{
				helper.FlushAllEvents(),
			},
		},
	}

	grill.Run(t, tests)
}
