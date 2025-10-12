package infrastructure

import (
	"context"
	"fmt"
	"log"
	"time"
)

// RetryConfig contains retry configuration
type RetryConfig struct {
	MaxRetries    int
	InitialDelay  time.Duration
	MaxDelay      time.Duration
	BackoffFactor float64
	EnableJitter  bool
}

// DefaultRetryConfig returns default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:    3,
		InitialDelay:  100 * time.Millisecond,
		MaxDelay:      5 * time.Second,
		BackoffFactor: 2.0,
		EnableJitter:  true,
	}
}

// RetryWithBackoff executes a function with exponential backoff retry
func RetryWithBackoff(ctx context.Context, config RetryConfig, operation func() error, operationName string) error {
	var lastErr error
	delay := config.InitialDelay

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		// First attempt (no delay)
		if attempt > 0 {
			// Apply jitter if enabled (random ±25% of delay)
			actualDelay := delay
			if config.EnableJitter {
				jitterFactor := 2.0*float64(time.Now().UnixNano()%100)/100.0 - 0.5
				jitter := time.Duration(float64(delay) * 0.25 * jitterFactor)
				actualDelay = delay + jitter
			}

			log.Printf("[Retry] %s: Attempt %d/%d failed. Retrying in %v... (error: %v)",
				operationName, attempt, config.MaxRetries, actualDelay, lastErr)

			select {
			case <-time.After(actualDelay):
				// Continue to retry
			case <-ctx.Done():
				return fmt.Errorf("retry cancelled: %w", ctx.Err())
			}
		}

		// Execute operation
		err := operation()
		if err == nil {
			if attempt > 0 {
				log.Printf("[Retry] %s: Succeeded on attempt %d/%d", operationName, attempt+1, config.MaxRetries+1)
			}
			return nil
		}

		lastErr = err

		// Calculate next delay with exponential backoff
		delay = time.Duration(float64(delay) * config.BackoffFactor)
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}
	}

	return fmt.Errorf("[Retry] %s: failed after %d attempts: %w", operationName, config.MaxRetries+1, lastErr)
}
