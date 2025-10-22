package messaging

import (
	"context"
	"log"
	"sync"

	"github.com/IBM/sarama"
	domainMq "github.com/kiin21/go-rest/services/notification-service/internal/notification/domain/messaging"
)

type KafkaNotificationConsumer struct {
	consumerGroup sarama.ConsumerGroup
	topics        []string
	handler       sarama.ConsumerGroupHandler
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

func NewKafkaNotificationConsumer(
	consumerGroup sarama.ConsumerGroup,
	topics []string,
	handler sarama.ConsumerGroupHandler,
) domainMq.NotificationConsumer {
	ctx, cancel := context.WithCancel(context.Background())
	return &KafkaNotificationConsumer{
		consumerGroup: consumerGroup,
		topics:        topics,
		handler:       handler,
		ctx:           ctx,
		cancel:        cancel,
	}
}

func (c *KafkaNotificationConsumer) Start() {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			if err := c.consumerGroup.Consume(c.ctx, c.topics, c.handler); err != nil {
				log.Printf("Error from consumer: %v", err)
			}

			// Check if context was cancelled
			if c.ctx.Err() != nil {
				log.Println("Consumer context cancelled, stopping...")
				return
			}
		}
	}()
	log.Printf("Kafka consumer started for topics: %v", c.topics)
}

func (c *KafkaNotificationConsumer) Stop() {
	log.Println("Stopping Kafka consumer...")
	c.cancel()  // Cancel context to stop consumption loop
	c.wg.Wait() // Wait for goroutine to finish

	if err := c.consumerGroup.Close(); err != nil {
		log.Printf("Error closing consumer group: %v", err)
	}
	log.Println("Kafka consumer stopped successfully")
}
