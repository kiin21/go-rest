package kafka

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"

	"github.com/kiin21/go-rest/pkg/events"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewProducer(brokers []string) (*Producer, error) {
	return newProducer(brokers, "")
}

func NewProducerWithTopic(brokers []string, topic string) (*Producer, error) {
	return newProducer(brokers, topic)
}

func newProducer(brokers []string, topic string) (*Producer, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka producer requires at least one broker")
	}

	config := defaultProducerConfig()

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
		topic:    topic,
	}, nil
}

func defaultProducerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.MaxMessageBytes = 1000000
	return config
}

func (p *Producer) PublishEvent(event *events.Event) (int32, int64, error) {
	return p.PublishEventToTopic("", event)
}

func (p *Producer) PublishEventToTopic(topic string, event *events.Event) (int32, int64, error) {
	if event == nil {
		return 0, 0, errors.New("event must not be nil")
	}
	if event.Type == "" {
		return 0, 0, errors.New("event type must not be empty")
	}

	if topic == "" {
		topic = p.topic
	}
	if topic == "" {
		return 0, 0, errors.New("kafka producer requires a topic")
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return 0, 0, fmt.Errorf("marshal event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.ByteEncoder(payload),
		Timestamp: event.Timestamp,
	}

	if event.Key != "" {
		msg.Key = sarama.StringEncoder(event.Key)
	} else if event.Domain != "" {
		msg.Key = sarama.StringEncoder(event.Domain)
	}

	partition, offset, err := p.SendProducerMessage(msg)
	if err != nil {
		return 0, 0, err
	}

	log.Printf("[KafkaProducer] Event sent - Topic: %s, Partition: %d, Offset: %d, Type: %s, Domain: %s",
		msg.Topic, partition, offset, event.Type, event.Domain)

	return partition, offset, nil
}

func (p *Producer) SendProducerMessage(msg *sarama.ProducerMessage) (int32, int64, error) {
	if msg == nil {
		return 0, 0, errors.New("kafka producer message is nil")
	}

	if msg.Topic == "" {
		msg.Topic = p.topic
	}
	if msg.Topic == "" {
		return 0, 0, errors.New("kafka producer requires a topic")
	}

	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return 0, 0, err
	}
	return partition, offset, nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
