package kafka

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/IBM/sarama"

	messaging "github.com/kiin21/go-rest/internal/shared/infrastructure/messagebroker"
)

// Consumer wraps a Kafka consumer group with retry support.
type Consumer struct {
	consumerGroup sarama.ConsumerGroup
	topic         string
	handler       func(context.Context, *messaging.SyncEvent) error
	ready         chan bool
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	retryConfig   messaging.RetryConfig
}

// NewConsumer creates a Kafka consumer group bound to the provided topic.
func NewConsumer(
	brokers []string,
	groupID string,
	topic string,
	handler func(context.Context, *messaging.SyncEvent) error,
) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRoundRobin}
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	log.Printf("[KafkaConsumer] Created consumer group: %s, topic: %s", groupID, topic)

	return &Consumer{
		consumerGroup: consumerGroup,
		topic:         topic,
		handler:       handler,
		ready:         make(chan bool),
		ctx:           ctx,
		cancel:        cancel,
		retryConfig:   messaging.DefaultRetryConfig(),
	}, nil
}

// Start begins consuming messages in a background goroutine.
func (kc *Consumer) Start() {
	kc.wg.Add(1)
	go func() {
		defer kc.wg.Done()
		for {
			if err := kc.consumerGroup.Consume(kc.ctx, []string{kc.topic}, kc); err != nil {
				log.Printf("[KafkaConsumer] Error from consumer: %v", err)
			}
			if kc.ctx.Err() != nil {
				return
			}
			kc.ready = make(chan bool)
		}
	}()

	<-kc.ready
	log.Printf("[KafkaConsumer] Consumer started and ready")
}

// Stop gracefully stops the consumer group.
func (kc *Consumer) Stop() {
	log.Printf("[KafkaConsumer] Stopping consumer...")
	kc.cancel()
	kc.wg.Wait()
	if err := kc.consumerGroup.Close(); err != nil {
		log.Printf("[KafkaConsumer] Error closing consumer: %v", err)
	}
	log.Printf("[KafkaConsumer] Consumer stopped")
}

// Setup runs at the beginning of a session.
func (kc *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(kc.ready)
	return nil
}

// Cleanup runs at the end of a session.
func (kc *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim processes messages from the claim.
func (kc *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			log.Printf("[KafkaConsumer] Message claimed: Topic=%s, Partition=%d, Offset=%d",
				message.Topic, message.Partition, message.Offset)

			var event messaging.SyncEvent
			if err := json.Unmarshal(message.Value, &event); err != nil {
				log.Printf("[KafkaConsumer] Failed to unmarshal message: %v", err)
				session.MarkMessage(message, "")
				continue
			}

			err := messaging.RetryWithBackoff(kc.ctx, kc.retryConfig, func() error {
				return kc.handler(kc.ctx, &event)
			}, event.Type)

			if err != nil {
				log.Printf("[KafkaConsumer] Failed to process event after retries: %v", err)
			} else {
				log.Printf("[KafkaConsumer] Event processed successfully: %s - %s", event.Type, event.Domain)
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}
