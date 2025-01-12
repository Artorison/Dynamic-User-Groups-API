package kafka

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(brokers []string) (*Producer, error) {
	config := GetKafkaConfig()

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{producer: producer}, nil
}

func (p *Producer) SendMessage(topic, key string, data interface{}) error {

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   (sarama.StringEncoder)(key),
		Value: (sarama.ByteEncoder)(jsonData),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Message sent to topic %s, partition %d, offset %d", topic, partition, offset)
	return nil
}

func (p *Producer) Close() {
	if err := p.producer.Close(); err != nil {
		log.Printf("Failed to close producer: %v", err)
	}
}
