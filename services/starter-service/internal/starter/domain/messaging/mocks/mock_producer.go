package mocks

import (
	"github.com/kiin21/go-rest/pkg/events"
)

// MockNotificationProducer is a mock implementation of NotificationProducer
type MockNotificationProducer struct {
	SendNotificationFunc func(event *events.Event) error
	CloseFunc            func() error
}

func (m *MockNotificationProducer) SendNotification(event *events.Event) error {
	if m.SendNotificationFunc != nil {
		return m.SendNotificationFunc(event)
	}
	return nil
}

func (m *MockNotificationProducer) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

