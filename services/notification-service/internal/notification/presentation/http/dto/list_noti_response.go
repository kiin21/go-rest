package dto

import (
	"time"

	"github.com/kiin21/go-rest/services/notification-service/internal/notification/domain/model"
)

type ListNotiResponse struct {
	ID          string    `json:"id"`
	FromStarter string    `json:"from_starter"`
	ToStarter   string    `json:"to_starter"`
	Message     string    `json:"message"`
	Type        string    `json:"type"`
	Timestamp   time.Time `json:"timestamp"`
}

func FromDomain(notifications []*model.Notification) []*ListNotiResponse {
	responses := make([]*ListNotiResponse, 0, len(notifications))
	for _, notification := range notifications {
		if notification == nil {
			continue
		}

		responses = append(responses, &ListNotiResponse{
			ID:          notification.ID,
			FromStarter: notification.FromStarter,
			ToStarter:   notification.ToStarter,
			Message:     notification.Message,
			Type:        notification.Type,
			Timestamp:   notification.Timestamp,
		})
	}

	return responses
}
