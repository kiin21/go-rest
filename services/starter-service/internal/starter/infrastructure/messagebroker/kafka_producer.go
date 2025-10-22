package messagebroker

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/kiin21/go-rest/pkg/events"
	domainMq "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/messaging"
)

type KafkaNotificationProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaNotificationProducer(producer sarama.SyncProducer, topic string) domainMq.NotificationProducer {
	return &KafkaNotificationProducer{
		producer: producer,
		topic:    topic,
	}
}

func (p *KafkaNotificationProducer) SendNotification(event *events.Event) error {
	// Marshal event to JSON
	eventBytes, err := event.ToBytes()
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(event.Key),
		Value: sarama.ByteEncoder(eventBytes),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		return err
	}

	log.Printf("Message sent to topic=%s, partition=%d, offset=%d", p.topic, partition, offset)
	return nil
}

func (p *KafkaNotificationProducer) Close() error {
	if p.producer != nil {
		return p.producer.Close()
	}
	return nil
}
