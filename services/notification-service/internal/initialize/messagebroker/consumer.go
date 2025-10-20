package messagebroker

import (
	"log"
	"strings"

	domainmessaging "github.com/kiin21/go-rest/services/notification-service/internal/notification/domain/messaging"
	inframessagebroker "github.com/kiin21/go-rest/services/notification-service/internal/notification/infrastructure/messagebroker"
)

func InitNotificationConsumer(
	kafkaBrokers,
	kafkaConsumerGroup,
	kafkaNotificationTopic string,
	service inframessagebroker.NotificationService,
) domainmessaging.NotificationConsumer {
	if service == nil {
		log.Printf("Warning: Kafka notification consumer service is nil, consumer disabled")
		return nil
	}

	brokers := parseKafkaBrokers(kafkaBrokers)
	if len(brokers) == 0 || kafkaNotificationTopic == "" || kafkaConsumerGroup == "" {
		log.Printf("Warning: Kafka notification consumer disabled due to incomplete configuration")
		return nil
	}

	consumer, err := inframessagebroker.NewKafkaNotificationConsumer(
		brokers,
		kafkaConsumerGroup,
		kafkaNotificationTopic,
		service,
	)
	if err != nil {
		log.Printf("Warning: failed to initialize Kafka notification consumer: %v", err)
		return nil
	}

	log.Printf("Kafka notification consumer initialised for topic %s", kafkaNotificationTopic)
	return consumer
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
