package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"project-survey-proceeder/contracts"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func InitProducer(url string) (contracts.IMessageProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	brokerList := []string{url}
	kafkaProducer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	producer := &KafkaProducer{kafkaProducer}
	return producer, nil
}

func (k *KafkaProducer) SendMessage(message []byte) error {
	topic := "my-topic"
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	_, _, err := k.producer.SendMessage(msg)
	if err != nil {
		fmt.Printf("Error producing message to Kafka: %v\n", err)
		return err
	}

	return nil
}

func (k *KafkaProducer) CloseConnection() error {
	return k.producer.Close()
}
