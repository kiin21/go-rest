package messagebroker

import (
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/kiin21/go-rest/pkg/utils"
	"github.com/kiin21/go-rest/services/starter-service/internal/config"
	domainMq "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/messaging"
	infraMq "github.com/kiin21/go-rest/services/starter-service/internal/starter/infrastructure/messagebroker"
)

func InitProducer(cfg config.Config) domainMq.NotificationProducer {
	brokers := utils.ParseString(cfg.KafkaBrokers, ",")
	if len(brokers) == 0 || cfg.KafkaTopicNotifications == "" {
		log.Printf("Warning: Kafka notification configuration incomplete, leader notifications disabled")
		return nil
	}

	// Create Sarama producer config
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll // Wait for all replicas
	saramaConfig.Producer.Retry.Max = 3
	saramaConfig.Producer.Compression = sarama.CompressionSnappy
	saramaConfig.Producer.Idempotent = true
	saramaConfig.Net.MaxOpenRequests = 1 // Required for idempotent producer
	saramaConfig.Producer.Retry.Backoff = 100 * time.Millisecond
	saramaConfig.Version = sarama.V2_8_0_0 // Kafka version
	saramaConfig.ClientID = "starter-service-producer"

	// Create sync producer
	producer, err := sarama.NewSyncProducer(brokers, saramaConfig)
	if err != nil {
		log.Printf("Warning: failed to initialize Kafka notification producer: %v", err)
		return nil
	}

	log.Printf("✅ Kafka notification producer initialized for topic: %s", cfg.KafkaTopicNotifications)

	return infraMq.NewKafkaNotificationProducer(producer, cfg.KafkaTopicNotifications)
}
