package document

import (
	"time"

	domainModel "github.com/kiin21/go-rest/services/notification-service/internal/notification/domain/model"
)

type NotificationDocument struct {
	ID          string    `bson:"_id"`
	FromStarter string    `bson:"from_starter"`
	ToStarter   string    `bson:"to_starter"`
	Message     string    `bson:"message"`
	Type        string    `bson:"type"`
	Timestamp   time.Time `bson:"timestamp"`
}

func (d *NotificationDocument) ToDomain() *domainModel.Notification {
	return &domainModel.Notification{
		ID:          d.ID,
		FromStarter: d.FromStarter,
		ToStarter:   d.ToStarter,
		Message:     d.Message,
		Type:        d.Type,
		Timestamp:   d.Timestamp,
	}
}
