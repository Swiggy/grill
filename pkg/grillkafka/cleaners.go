package grillkafka

import (
	"context"
	"time"

	"bitbucket.org/swigy/grill"
	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
)

func (grillkafka *GrillKafka) DeleteTopics(topics ...string) grill.Cleaner {
	return grill.CleanerFunc(func() error {
		_, err := grillkafka.kafka.AdminClient.DeleteTopics(context.TODO(), topics, confluent.SetAdminOperationTimeout(time.Second*10))
		return err
	})
}
