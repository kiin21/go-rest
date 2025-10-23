package messagebroker

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/kiin21/go-rest/pkg/utils"
	"github.com/kiin21/go-rest/services/starter-service/internal/config"
	domainMq "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/messaging"
	domainRepo "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository"
	domainService "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/service"
	infraMq "github.com/kiin21/go-rest/services/starter-service/internal/starter/infrastructure/messagebroker"
)

func InitEventHandler(
	starterSearchRepo domainRepo.StarterSearchRepository,
	starterRepo domainRepo.StarterRepository,
	enrichmentService *domainService.StarterEnrichmentService,
) *EventHandler {
	return NewEventHandler(starterRepo, starterSearchRepo, enrichmentService)
}

func InitGroupConsumer(cfg config.Config, handler *EventHandler) domainMq.StarterConsumer {
	// 1. Setup Sarama config
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	saramaConfig.Version = sarama.V2_8_0_0

	// 2. Create a consumer group
	consumerGroup, err := sarama.NewConsumerGroup(
		utils.ParseString(cfg.KafkaBrokers, ","),
		cfg.KafkaConsumerGroup,
		saramaConfig,
	)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}

	// 3. Create a SyncConsumer - listen to sync events topic for Elasticsearch indexing
	topics := []string{cfg.KafkaTopicSyncEvents}
	consumer := infraMq.NewKafkaStarterConsumer(consumerGroup, topics, handler)

	// 4. Start consuming
	consumer.Start()

	return consumer
}
