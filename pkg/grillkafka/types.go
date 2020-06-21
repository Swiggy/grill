package grillkafka

import (
	confluent "github.com/confluentinc/confluent-kafka-go/kafka"
)

type Message struct {
	Topic   string
	Key     string
	Value   string
	Headers map[string]string
}

type Consumer struct {
	msgChan       chan Message
	kafkaConsumer *confluent.Consumer
}

func (c *Consumer) Close() {
	if isOpen(c.msgChan) {
		close(c.msgChan)
	}
	go func() {
		c.kafkaConsumer.Unsubscribe()
		c.kafkaConsumer.Close()
	}()
}

func (c *Consumer) MsgChan() <-chan Message {
	return c.msgChan
}

func isOpen(ch <-chan Message) bool {
	select {
	case <-ch:
		return false
	default:
		return true
	}
}
