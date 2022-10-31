package grillkafka

import (
	"context"
	"github.com/Swiggy/grill/canned"
	"time"

	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
)

type Kafka struct {
	kafka *canned.Kafka
}

func (gk *Kafka) Start(ctx context.Context) error {
	kafka, err := canned.NewKafka(ctx)
	if err != nil {
		return err
	}
	gk.kafka = kafka
	return nil
}

func (gk *Kafka) BootstrapServers() string {
	return gk.kafka.BootstrapServers
}

func (gk *Kafka) Stop(ctx context.Context) error {
	return gk.kafka.Container.Terminate(ctx)
}

func (gk *Kafka) Produce(topic string, message Message) error {
	deliveryChan := make(chan confluent.Event)
	var headers []confluent.Header
	for key, value := range message.Headers {
		headers = append(headers, confluent.Header{Key: key, Value: []byte(value)})
	}
	err := gk.kafka.Producer.Produce(&confluent.Message{
		TopicPartition: confluent.TopicPartition{Topic: &topic, Partition: confluent.PartitionAny},
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

func (gk *Kafka) NewConsumer(group, topic string, openTime time.Duration) (*Consumer, error) {
	kafkaConsumer, err := confluent.NewConsumer(&confluent.ConfigMap{
		"bootstrap.servers":               gk.kafka.BootstrapServers,
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
