package grillkafka

import (
	"context"
	"time"

	"bitbucket.org/swigy/grill"
	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
)

func (grillkafka *GrillKafka) CreateTopics(topics ...string) grill.Stub {
	return grill.StubFunc(func() error {
		var ts []confluent.TopicSpecification
		for _, topic := range topics {
			ts = append(ts, confluent.TopicSpecification{
				Topic:             topic,
				NumPartitions:     1,
				ReplicationFactor: 1,
			})
		}
		_, err := grillkafka.kafka.AdminClient.CreateTopics(context.TODO(), ts, confluent.SetAdminOperationTimeout(time.Second*10))
		return err
	})
}
