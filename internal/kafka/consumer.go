package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumer sarama.Consumer
}

func NewConsumer(brokers []string) (*Consumer, error) {
	config := GetKafkaConfig()
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{consumer: consumer}, nil
}

func (c *Consumer) ConsumeMessages(topic string, processFunc func(message *sarama.ConsumerMessage) error) {
	partitionConsumer, err := c.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to start consumer for topic %s: %v", topic, err)
	}
	defer partitionConsumer.Close()

	log.Printf("Started consumer for topic: %s", topic)

	for msg := range partitionConsumer.Messages() {
		log.Printf("Message received from topic %s: %s", topic, string(msg.Value))

		if err := processFunc(msg); err != nil {
			log.Printf("Failed to process message: %v", err)
		}
	}
}

func (c *Consumer) Close() {
	if err := c.consumer.Close(); err != nil {
		log.Printf("Failed to close consumer: %v", err)
	}
}
