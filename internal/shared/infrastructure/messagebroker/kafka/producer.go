package kafka

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/kiin21/go-rest/internal/shared/infrastructure/messagebroker"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.MaxMessageBytes = 1000000

	p, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	log.Printf("[KafkaProducer] Connected to Kafka brokers: %v, topic: %s", brokers, topic)

	return &Producer{
		producer: p,
		topic:    topic,
	}, nil
}

func (kp *Producer) PublishSyncEvent(event *messaging.SyncEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("[KafkaProducer] Failed to marshal event: %v", err)
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic:     kp.topic,
		Key:       sarama.StringEncoder(event.Domain),
		Value:     sarama.ByteEncoder(payload),
		Timestamp: time.Now(),
	}

	partition, offset, err := kp.producer.SendMessage(msg)
	if err != nil {
		log.Printf("[KafkaProducer] Failed to send message: %v", err)
		return err
	}

	log.Printf("[KafkaProducer] Message sent successfully - Topic: %s, Partition: %d, Offset: %d, Domain: %s, Type: %s",
		kp.topic, partition, offset, event.Domain, event.Type)

	return nil
}

func (kp *Producer) Close() error {
	log.Printf("[KafkaProducer] Closing producer...")
	return kp.producer.Close()
}
