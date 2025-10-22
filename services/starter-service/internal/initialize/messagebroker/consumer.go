package messagebroker

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/kiin21/go-rest/pkg/utils"
	"github.com/kiin21/go-rest/services/starter-service/internal/config"
	domainMq "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/messaging"
	domainRepo "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository"
	infraMq "github.com/kiin21/go-rest/services/starter-service/internal/starter/infrastructure/messagebroker"
)

func InitEventHandler(repo domainRepo.StarterSearchRepository) *EventHandler {
	return NewEventHandler(repo)
}

func InitGroupConsumer(cfg config.Config, handler *EventHandler) domainMq.StarterConsumer {
	// 1. Setup Sarama config
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	saramaConfig.Version = sarama.V2_8_0_0

	// 2. Create a consumer group
	consumerGroup, err := sarama.NewConsumerGroup(
		utils.ParseCSVString(cfg.KafkaBrokers),
		cfg.KafkaConsumerGroup,
		saramaConfig,
	)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}

	// 3. Create a NotificationConsumer
	topics := []string{cfg.KafkaTopicNotifications}
	consumer := infraMq.NewKafkaStarterConsumer(consumerGroup, topics, handler)

	// 4. Start consuming
	consumer.Start()

	return consumer
}
