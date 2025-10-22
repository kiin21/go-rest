package messagequeue

import "time"

type ConsumerConfig struct {
	Brokers           []string
	GroupID           string
	Topics            []string
	Version           string
	SessionTimeout    time.Duration
	HeartbeatInterval time.Duration
	InitialOffset     string // "oldest" or "newest"
	FetchDefaultBytes int32
	RebalanceStrategy string // "range", "roundrobin", "sticky"
}

func DefaultConsumerConfig() *ConsumerConfig {
	return &ConsumerConfig{
		Brokers:           []string{"localhost:9092"},
		Version:           "2.8.0",
		SessionTimeout:    30 * time.Second,
		HeartbeatInterval: 3 * time.Second,
		InitialOffset:     "oldest",
		FetchDefaultBytes: 8 * 1024,
		RebalanceStrategy: "range",
	}
}

type ProducerConfig struct {
	Brokers          []string
	Version          string
	ClientID         string
	CompressionType  string
	MaxRetry         int
	RetryBackoff     time.Duration
	RequiredAcks     int
	EnableIdempotent bool
}

func DefaultProducerConfig() *ProducerConfig {
	return &ProducerConfig{
		Brokers:          []string{"localhost:9092"},
		Version:          "2.8.0",
		ClientID:         "default-producer",
		CompressionType:  "snappy",
		MaxRetry:         3,
		RetryBackoff:     100 * time.Millisecond,
		RequiredAcks:     -1, // WaitForAll
		EnableIdempotent: true,
	}
}
