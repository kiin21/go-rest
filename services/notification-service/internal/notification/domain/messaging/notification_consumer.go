package messaging

type NotificationConsumer interface {
	Start()
	Stop()
}
