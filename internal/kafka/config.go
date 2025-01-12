package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

func GetKafkaConfig() *sarama.Config {
	config := sarama.NewConfig()

	kafkaVersion, err := sarama.ParseKafkaVersion("2.8.0")
	if err != nil {
		log.Fatalf("Error parsing Kafka version: %v", err)
	}
	config.Version = kafkaVersion

	config.Producer.Return.Successes = true
	config.Consumer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Partitioner = sarama.NewManualPartitioner

	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
