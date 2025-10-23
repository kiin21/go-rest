package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kiin21/go-rest/pkg/mongodb"
	"github.com/kiin21/go-rest/services/notification-service/internal/config"
	"go.mongodb.org/mongo-driver/event"
)

func InitDB(cfg *config.Config) (*mongodb.Client, error) {
	if cfg.MongoURI == "" {
		return nil, fmt.Errorf("mongodb uri is not configured")
	}
	if cfg.MongoDatabase == "" {
		return nil, fmt.Errorf("mongodb database is not configured")
	}

	logger := log.New(os.Stdout, "[MongoDB] ", log.LstdFlags)

	cmdMonitor := &event.CommandMonitor{
		Started: func(ctx context.Context, e *event.CommandStartedEvent) {
			logger.Printf("[QUERY] Database: %s | Command: %s | Request ID: %d",
				e.DatabaseName,
				e.CommandName,
				e.RequestID,
			)
		},
		Succeeded: func(ctx context.Context, e *event.CommandSucceededEvent) {
			// Highlight slow queries
			if e.Duration > time.Second {
				logger.Printf("[SLOW QUERY] Command: %s | Duration: %v | Request ID: %d",
					e.CommandName,
					e.Duration,
					e.RequestID,
				)
			} else {
				logger.Printf("[SUCCESS] Command: %s | Duration: %v | Request ID: %d",
					e.CommandName,
					e.Duration,
					e.RequestID,
				)
			}
		},
		Failed: func(ctx context.Context, e *event.CommandFailedEvent) {
			logger.Printf("[ERROR] Command: %s | Duration: %v | Failure: %s | Request ID: %d",
				e.CommandName,
				e.Duration,
				e.Failure,
				e.RequestID,
			)
		},
	}

	poolMonitor := &event.PoolMonitor{
		Event: func(e *event.PoolEvent) {
			switch e.Type {
			case event.ConnectionCreated:
				logger.Printf("[POOL] Connection created | Address: %s", e.Address)
			case event.ConnectionClosed:
				logger.Printf("[POOL] Connection closed | Address: %s | Reason: %s", e.Address, e.Reason)
			case event.GetSucceeded:
				logger.Printf("[POOL] Connection checked out | Address: %s | Duration: %v", e.Address, e.Duration)
			case event.GetFailed:
				logger.Printf("[POOL ERROR] Failed to get connection | Address: %s | Reason: %s", e.Address, e.Reason)
			case event.ConnectionReady:
				logger.Printf("[POOL] Connection ready | Address: %s", e.Address)
			case event.PoolCleared:
				logger.Printf("[POOL] Pool cleared | Address: %s", e.Address)
			}
		},
	}

	client, err := mongodb.NewClient(
		cfg.MongoURI,
		cfg.MongoDatabase,
		cmdMonitor,
		poolMonitor,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize mongodb client: %w", err)
	}

	log.Println("MongoDB connection established successfully.")
	return client, nil
}
