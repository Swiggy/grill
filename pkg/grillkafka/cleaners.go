package grillkafka

import (
	"context"
	"time"

	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/lovlin-thakkar/swiggy-grill"
)

func (gk *Kafka) DeleteTopics(topics ...string) grill.Cleaner {
	return grill.CleanerFunc(func() error {
		_, err := gk.kafka.AdminClient.DeleteTopics(context.TODO(), topics, confluent.SetAdminOperationTimeout(time.Second*10))
		return err
	})
}
