package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"

	"github.com/IBM/sarama"

	"github.com/kiin21/go-rest/pkg/events"
)

type Consumer struct {
	consumer sarama.ConsumerGroup
}

type MessageHandler func(message []byte) error

func NewConsumer(brokers []string, groupID string) (*Consumer, error) {
	group, err := newConsumerGroup(brokers, groupID)
	if err != nil {
		return nil, err
	}
	return &Consumer{consumer: group}, nil
}

func newConsumerGroup(brokers []string, groupID string) (sarama.ConsumerGroup, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka consumer requires at least one broker")
	}
	if groupID == "" {
		return nil, errors.New("kafka consumer requires a group id")
	}

	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true

	return sarama.NewConsumerGroup(brokers, groupID, config)
}

func (c *Consumer) ConsumerGroup() sarama.ConsumerGroup {
	return c.consumer
}

func (c *Consumer) Subscribe(ctx context.Context, topics []string, handler MessageHandler) error {
	if len(topics) == 0 {
		return errors.New("at least one topic is required")
	}
	if handler == nil {
		return errors.New("message handler is required")
	}

	consumerHandler := &consumerGroupHandler{handler: handler}

	for {
		if err := c.consumer.Consume(ctx, topics, consumerHandler); err != nil {
			log.Printf("[KafkaConsumer] Error from consumer: %v", err)
			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func (c *Consumer) Close() error {
	if c == nil || c.consumer == nil {
		return nil
	}
	return c.consumer.Close()
}

type consumerGroupHandler struct {
	handler MessageHandler
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			if err := h.handler(message.Value); err != nil {
				log.Printf("[KafkaConsumer] Error handling message: %v", err)
			} else {
				session.MarkMessage(message, "")
			}

		case <-session.Context().Done():
			return nil
		}
	}
}

type EventHandler func(context.Context, *events.Event) error

type EventConsumer struct {
	consumerGroup sarama.ConsumerGroup
	handler       EventHandler
	topics        []string
	retryConfig   RetryConfig

	ready  chan struct{}
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewEventConsumer(
	brokers []string,
	groupID string,
	topics []string,
	handler EventHandler,
	opts ...EventConsumerOption,
) (*EventConsumer, error) {
	if len(topics) == 0 {
		return nil, errors.New("event consumer requires at least one topic")
	}
	if handler == nil {
		return nil, errors.New("event consumer handler must not be nil")
	}

	group, err := newConsumerGroup(brokers, groupID)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	consumer := &EventConsumer{
		consumerGroup: group,
		handler:       handler,
		topics:        topics,
		retryConfig:   DefaultRetryConfig(),
		ready:         make(chan struct{}),
		ctx:           ctx,
		cancel:        cancel,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(consumer)
		}
	}

	return consumer, nil
}

type EventConsumerOption func(*EventConsumer)

func WithEventRetryConfig(cfg RetryConfig) EventConsumerOption {
	return func(c *EventConsumer) {
		c.retryConfig = cfg
	}
}

func WithEventContext(ctx context.Context) EventConsumerOption {
	return func(c *EventConsumer) {
		if ctx == nil {
			return
		}
		c.ctx, c.cancel = context.WithCancel(ctx)
	}
}

func (c *EventConsumer) Start() {
	if c == nil {
		return
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			if err := c.consumerGroup.Consume(c.ctx, c.topics, c); err != nil {
				log.Printf("[KafkaEventConsumer] Consume error: %v", err)
			}
			if c.ctx.Err() != nil {
				return
			}
			c.ready = make(chan struct{})
		}
	}()

	<-c.ready
	log.Printf("[KafkaEventConsumer] Consumer started for topics: %v", c.topics)
}

func (c *EventConsumer) Stop() {
	if c == nil {
		return
	}
	log.Printf("[KafkaEventConsumer] Stopping consumer...")
	c.cancel()
	c.wg.Wait()
	if err := c.consumerGroup.Close(); err != nil {
		log.Printf("[KafkaEventConsumer] Error closing consumer: %v", err)
	}
}

func (c *EventConsumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

func (c *EventConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *EventConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			var event events.Event
			if err := json.Unmarshal(message.Value, &event); err != nil {
				log.Printf("[KafkaEventConsumer] Failed to unmarshal event: %v", err)
				session.MarkMessage(message, "")
				continue
			}

			handlerCtx := session.Context()

			err := RetryWithBackoff(handlerCtx, c.retryConfig, func() error {
				return c.handler(handlerCtx, &event)
			}, event.Type)

			if err != nil {
				log.Printf("[KafkaEventConsumer] Handler failed after retries: %v", err)
			} else {
				log.Printf("[KafkaEventConsumer] Event processed successfully: %s - %s", event.Type, event.Domain)
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}
