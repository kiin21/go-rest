package messagebroker

import (
	"context"
	"log"
	"sync"

	"github.com/IBM/sarama"
	domainMq "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/messaging"
)

type KafkaStarterConsumer struct {
	consumerGroup sarama.ConsumerGroup
	topics        []string
	handler       sarama.ConsumerGroupHandler
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

func NewKafkaStarterConsumer(
	consumerGroup sarama.ConsumerGroup,
	topics []string,
	handler sarama.ConsumerGroupHandler,
) domainMq.StarterConsumer {
	ctx, cancel := context.WithCancel(context.Background())
	return &KafkaStarterConsumer{
		consumerGroup: consumerGroup,
		topics:        topics,
		handler:       handler,
		ctx:           ctx,
		cancel:        cancel,
	}
}

func (c *KafkaStarterConsumer) Start() {
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

func (c *KafkaStarterConsumer) Stop() {
	log.Println("Stopping Kafka consumer...")
	c.cancel()
	c.wg.Wait()

	if err := c.consumerGroup.Close(); err != nil {
		log.Printf("Error closing consumer group: %v", err)
	}
	log.Println("Kafka consumer stopped successfully")
}
