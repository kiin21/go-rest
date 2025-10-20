package messaging

import (
	"context"

	"github.com/kiin21/go-rest/pkg/events"
)

type NotificationPublisher interface {
	PublishLeaderAssignment(context.Context, *events.LeaderAssignmentNotification) error
}
