package mocks

import (
	"github.com/kiin21/go-rest/pkg/events"
)

// MockSyncProducer is a mock implementation of SyncProducer
type MockSyncProducer struct {
	SendSyncEventFunc func(event *events.Event) error
	CloseFunc         func() error
}

func (m *MockSyncProducer) SendSyncEvent(event *events.Event) error {
	if m.SendSyncEventFunc != nil {
		return m.SendSyncEventFunc(event)
	}
	return nil
}

func (m *MockSyncProducer) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

