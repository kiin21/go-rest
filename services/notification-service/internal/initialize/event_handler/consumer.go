package event_handler

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/kiin21/go-rest/pkg/utils"
	"github.com/kiin21/go-rest/services/notification-service/internal/config"
	domainMq "github.com/kiin21/go-rest/services/notification-service/internal/notification/domain/messaging"
	domainRepo "github.com/kiin21/go-rest/services/notification-service/internal/notification/domain/repository"
	infraMq "github.com/kiin21/go-rest/services/notification-service/internal/notification/infrastructure/messaging"
)

func InitEventHandler(repo domainRepo.NotificationRepository) *EventHandler {
	return NewEventHandler(repo)
}

func InitGroupConsumer(cfg config.Config, handler *EventHandler) domainMq.NotificationConsumer {
	// 1. Setup Sarama config
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	saramaConfig.Version = sarama.V2_8_0_0

	// 2. Create a consumer group
	consumerGroup, err := sarama.NewConsumerGroup(
		utils.ParseCSVString(cfg.KafkaBrokers, ","),
		cfg.KafkaConsumerGroup,
		saramaConfig,
	)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}

	// 3. Create a NotificationConsumer
	topics := []string{cfg.KafkaTopicNotifications}
	consumer := infraMq.NewKafkaNotificationConsumer(consumerGroup, topics, handler)

	// 4. Start consuming
	consumer.Start()

	return consumer
}
