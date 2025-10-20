package messagebroker

import (
	"log"
	"strings"

	sharedKafka "github.com/kiin21/go-rest/pkg/kafka"
	domainmessaging "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/messaging"
	broker "github.com/kiin21/go-rest/services/starter-service/internal/starter/infrastructure/messagebroker"
)

func InitNotificationPublisher(kafkaBrokers, kafkaNotificationTopic string) domainmessaging.NotificationPublisher {
	brokers := parseKafkaBrokers(kafkaBrokers)
	if len(brokers) == 0 || kafkaNotificationTopic == "" {
		log.Printf("Warning: Kafka notification configuration incomplete, leader notifications disabled")
		return nil
	}

	producer, err := sharedKafka.NewProducer(brokers)
	if err != nil {
		log.Printf("Warning: failed to initialize Kafka notification producer: %v", err)
		return nil
	}

	log.Printf("Kafka notification publisher initialised for topic %s", kafkaNotificationTopic)
	return broker.NewKafkaNotificationPublisher(producer, kafkaNotificationTopic)
}

func parseKafkaBrokers(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	brokers := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			brokers = append(brokers, trimmed)
		}
	}
	return brokers
}
