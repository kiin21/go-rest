package kafka

import (
	"context"
	"fmt"
	"log"
	"time"
)

type RetryConfig struct {
	MaxRetries    int
	InitialDelay  time.Duration
	MaxDelay      time.Duration
	BackoffFactor float64
	EnableJitter  bool
}

func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:    3,
		InitialDelay:  100 * time.Millisecond,
		MaxDelay:      5 * time.Second,
		BackoffFactor: 2.0,
		EnableJitter:  true,
	}
}

func RetryWithBackoff(ctx context.Context, config RetryConfig, operation func() error, operationName string) error {
	var lastErr error
	delay := config.InitialDelay

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		if attempt > 0 {
			actualDelay := delay
			if config.EnableJitter {
				jitterFactor := 2.0*float64(time.Now().UnixNano()%100)/100.0 - 0.5
				jitter := time.Duration(float64(delay) * 0.25 * jitterFactor)
				actualDelay = delay + jitter
			}

			log.Printf("[KafkaRetry] %s: attempt %d/%d failed. Retrying in %v... (error: %v)",
				operationName, attempt, config.MaxRetries, actualDelay, lastErr)

			select {
			case <-time.After(actualDelay):
			case <-ctx.Done():
				return fmt.Errorf("retry cancelled: %w", ctx.Err())
			}
		}

		if err := operation(); err == nil {
			if attempt > 0 {
				log.Printf("[KafkaRetry] %s: succeeded on attempt %d/%d", operationName, attempt+1, config.MaxRetries+1)
			}
			return nil
		} else {
			lastErr = err
			delay = time.Duration(float64(delay) * config.BackoffFactor)
			if delay > config.MaxDelay {
				delay = config.MaxDelay
			}
		}
	}

	return fmt.Errorf("[KafkaRetry] %s: failed after %d attempts: %w", operationName, config.MaxRetries+1, lastErr)
}
