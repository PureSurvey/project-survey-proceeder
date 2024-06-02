package kafka

import (
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"log"
	"os"
	"project-survey-proceeder/internal/configuration"
	"time"
)

type Producer struct {
	producer sarama.AsyncProducer
	config   *configuration.EventsConfiguration
}

func NewProducer(config *configuration.EventsConfiguration) *Producer {
	return &Producer{config: config}
}

func (p *Producer) Init() error {
	config := sarama.NewConfig()
	config.Metadata.Retry.Backoff = 500 * time.Millisecond

	brokerList := []string{p.config.BrokerUrl}
	kafkaProducer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		return err
	}

	p.producer = kafkaProducer
	go p.listenToErrors()

	return nil
}

func (p *Producer) AsyncSendMessage(message []byte) {
	msg := &sarama.ProducerMessage{
		Topic: p.config.Topic,
		Value: sarama.StringEncoder(message),
	}

	p.producer.Input() <- msg
}

func (p *Producer) CloseConnection() {
	if p.producer != nil {
		p.producer.AsyncClose()
	}
}

func (p *Producer) listenToErrors() {
	for {
		err := <-p.producer.Errors()
		log.Printf(err.Error())
		filename := time.Now().UTC().String() + "_" + uuid.New().String()
		val, _ := err.Msg.Value.Encode()
		os.WriteFile(filename, val, os.ModeDir)
	}
}
