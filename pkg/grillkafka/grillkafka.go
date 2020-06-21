package grillkafka

import (
	"context"
	"time"

	"bitbucket.org/swigy/grill/internal/canned"
	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
)

type GrillKafka struct {
	kafka *canned.Kafka
}

func Start() (*GrillKafka, error) {
	kafka, err := canned.NewKafka(context.TODO())
	if err != nil {
		return nil, err
	}
	return &GrillKafka{
		kafka: kafka,
	}, nil
}

func (grillkafka *GrillKafka) BootstrapServers() string {
	return grillkafka.kafka.BootstrapServers
}

func (grillkafka *GrillKafka) Stop() error {
	return grillkafka.kafka.Container.Terminate(context.Background())
}

func (grillkafka *GrillKafka) Produce(message Message) error {
	deliveryChan := make(chan confluent.Event)
	var headers []confluent.Header
	for key, value := range message.Headers {
		headers = append(headers, confluent.Header{Key: key, Value: []byte(value)})
	}
	err := grillkafka.kafka.Producer.Produce(&confluent.Message{
		TopicPartition: confluent.TopicPartition{Topic: &message.Topic, Partition: confluent.PartitionAny},
		Key:            []byte(message.Key),
		Value:          []byte(message.Value),
		Headers:        headers,
	}, deliveryChan)
	if err != nil {
		return err
	}
	e := <-deliveryChan
	m := e.(*confluent.Message)
	return m.TopicPartition.Error
}

func (grillkafka *GrillKafka) NewConsumer(group, topic string, openTime time.Duration) (*Consumer, error) {
	kafkaConsumer, err := confluent.NewConsumer(&confluent.ConfigMap{
		"bootstrap.servers":               grillkafka.kafka.BootstrapServers,
		"broker.address.family":           "v4",
		"group.id":                        group,
		"session.timeout.ms":              6000,
		"auto.offset.reset":               "earliest",
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
	})
	if err != nil {
		return nil, err
	}
	err = kafkaConsumer.SubscribeTopics([]string{topic}, nil)
	outputChan := make(chan Message, 100)
	consumer := &Consumer{msgChan: outputChan, kafkaConsumer: kafkaConsumer}

	go func() {
		for {
			select {
			case ev := <-kafkaConsumer.Events():
				switch e := ev.(type) {
				case confluent.AssignedPartitions:
					kafkaConsumer.Assign(e.Partitions)
				case confluent.RevokedPartitions:
					kafkaConsumer.Unassign()
				case *confluent.Message:
					headers := map[string]string{}
					for _, h := range e.Headers {
						headers[h.Key] = string(h.Value)
					}
					outputChan <- Message{
						Topic:   *e.TopicPartition.Topic,
						Key:     string(e.Key),
						Value:   string(e.Value),
						Headers: headers,
					}
				}
			case <-time.Tick(openTime):
				consumer.Close()
				return
			}
		}
	}()
	return consumer, nil
}
