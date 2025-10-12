package infrastructure

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

// KafkaConsumer handles consuming messages from Kafka
type KafkaConsumer struct {
	consumerGroup sarama.ConsumerGroup
	topic         string
	handler       func(context.Context, *SyncEvent) error
	ready         chan bool
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	retryConfig   RetryConfig
}

// NewKafkaConsumer creates a new Kafka consumer
func NewKafkaConsumer(
	brokers []string,
	groupID string,
	topic string,
	handler func(context.Context, *SyncEvent) error,
) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	log.Printf("[KafkaConsumer] Created consumer group: %s, topic: %s", groupID, topic)

	return &KafkaConsumer{
		consumerGroup: consumerGroup,
		topic:         topic,
		handler:       handler,
		ready:         make(chan bool),
		ctx:           ctx,
		cancel:        cancel,
		retryConfig:   DefaultRetryConfig(),
	}, nil
}

// Start starts consuming messages from Kafka
func (kc *KafkaConsumer) Start() {
	kc.wg.Add(1)
	go func() {
		defer kc.wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := kc.consumerGroup.Consume(kc.ctx, []string{kc.topic}, kc); err != nil {
				log.Printf("[KafkaConsumer] Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if kc.ctx.Err() != nil {
				return
			}
			kc.ready = make(chan bool)
		}
	}()

	<-kc.ready // Wait till the consumer has been set up
	log.Printf("[KafkaConsumer] Consumer started and ready")
}

// Stop stops the consumer
func (kc *KafkaConsumer) Stop() {
	log.Printf("[KafkaConsumer] Stopping consumer...")
	kc.cancel()
	kc.wg.Wait()
	if err := kc.consumerGroup.Close(); err != nil {
		log.Printf("[KafkaConsumer] Error closing consumer: %v", err)
	}
	log.Printf("[KafkaConsumer] Consumer stopped")
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (kc *KafkaConsumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(kc.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (kc *KafkaConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages()
func (kc *KafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			log.Printf("[KafkaConsumer] Message claimed: Topic=%s, Partition=%d, Offset=%d",
				message.Topic, message.Partition, message.Offset)

			// Deserialize message
			var event SyncEvent
			if err := json.Unmarshal(message.Value, &event); err != nil {
				log.Printf("[KafkaConsumer] Failed to unmarshal message: %v", err)
				session.MarkMessage(message, "")
				continue
			}

			// Process event with retry logic
			err := RetryWithBackoff(kc.ctx, kc.retryConfig, func() error {
				return kc.handler(kc.ctx, &event)
			}, event.Type)

			if err != nil {
				log.Printf("[KafkaConsumer] Failed to process event after retries: %v", err)
				// TODO: Send to dead letter queue (separate Kafka topic)
				// For now, we mark it as processed to avoid blocking
			} else {
				log.Printf("[KafkaConsumer] Event processed successfully: %s - %s", event.Type, event.Domain)
			}

			// Mark message as processed
			session.MarkMessage(message, "")

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}
