package model

import (
	"time"
)

type Notification struct {
	ID          string
	FromStarter string
	ToStarter   string
	Message     string
	Type        string
	Timestamp   time.Time
}
