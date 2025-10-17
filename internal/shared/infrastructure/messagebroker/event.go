package messaging

import "time"

type SyncEvent struct {
	Type      string
	Domain    string
	Data      interface{}
	Timestamp time.Time
	Retries   int
}
