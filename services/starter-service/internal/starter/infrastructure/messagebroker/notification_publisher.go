package messagebroker

import (
	"context"

	"github.com/kiin21/go-rest/pkg/events"
	pkgKafka "github.com/kiin21/go-rest/pkg/kafka"
	domainmessaging "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/messaging"
)

type kafkaNotificationPublisher struct {
	producer *pkgKafka.Producer
	topic    string
}

func NewKafkaNotificationPublisher(producer *pkgKafka.Producer, topic string) domainmessaging.NotificationPublisher {
	return &kafkaNotificationPublisher{
		producer: producer,
		topic:    topic,
	}
}

func (p *kafkaNotificationPublisher) PublishLeaderAssignment(ctx context.Context, event *events.LeaderAssignmentNotification) error {
	_ = ctx
	if p == nil || p.producer == nil || event == nil {
		return nil
	}
	kafkaEvent, err := events.NewLeaderAssignmentEvent(event)
	if err != nil {
		return err
	}
	_, _, err = p.producer.PublishEventToTopic(p.topic, kafkaEvent)
	return err
}
