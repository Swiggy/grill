package grillkafka

import (
	"context"
	"time"

	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/lovlin-thakkar/swiggy-grill"
)

func (gk *Kafka) CreateTopics(topics ...string) grill.Stub {
	return grill.StubFunc(func() error {
		var ts []confluent.TopicSpecification
		for _, topic := range topics {
			ts = append(ts, confluent.TopicSpecification{
				Topic:             topic,
				NumPartitions:     1,
				ReplicationFactor: 1,
			})
		}
		_, err := gk.kafka.AdminClient.CreateTopics(context.TODO(), ts, confluent.SetAdminOperationTimeout(time.Second*10))
		return err
	})
}
