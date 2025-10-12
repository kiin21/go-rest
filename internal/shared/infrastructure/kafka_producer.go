package infrastructure

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
)

// KafkaProducer handles sending messages to Kafka
type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer(brokers []string, topic string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.MaxMessageBytes = 1000000

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	log.Printf("[KafkaProducer] Connected to Kafka brokers: %v, topic: %s", brokers, topic)

	return &KafkaProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

// PublishSyncEvent publishes a sync event to Kafka
func (kp *KafkaProducer) PublishSyncEvent(ctx context.Context, event *SyncEvent) error {
	// Serialize event to JSON
	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("[KafkaProducer] Failed to marshal event: %v", err)
		return err
	}

	// Create Kafka message
	msg := &sarama.ProducerMessage{
		Topic:     kp.topic,
		Key:       sarama.StringEncoder(event.Domain), // Use domain as key for partitioning
		Value:     sarama.ByteEncoder(eventBytes),
		Timestamp: time.Now(),
	}

	// Send message
	partition, offset, err := kp.producer.SendMessage(msg)
	if err != nil {
		log.Printf("[KafkaProducer] Failed to send message: %v", err)
		return err
	}

	log.Printf("[KafkaProducer] Message sent successfully - Topic: %s, Partition: %d, Offset: %d, Domain: %s, Type: %s",
		kp.topic, partition, offset, event.Domain, event.Type)

	return nil
}

// Close closes the Kafka producer
func (kp *KafkaProducer) Close() error {
	log.Printf("[KafkaProducer] Closing producer...")
	return kp.producer.Close()
}
