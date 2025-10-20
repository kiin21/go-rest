package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	MongoClient *mongo.Client
	Database    *mongo.Database
}

func NewClient(
	uri string,
	database string,
	cmdMonitor *event.CommandMonitor,
	poolMonitor *event.PoolMonitor,
) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)

	clientOpts.SetMaxPoolSize(100)           // MaxOpenConns
	clientOpts.SetMinPoolSize(10)            // MaxIdleConns
	clientOpts.SetMaxConnIdleTime(time.Hour) // ConnMaxLifetime

	clientOpts.SetConnectTimeout(10 * time.Second)
	clientOpts.SetServerSelectionTimeout(5 * time.Second)

	if cmdMonitor != nil {
		clientOpts.SetMonitor(cmdMonitor)
	}
	if poolMonitor != nil {
		clientOpts.SetPoolMonitor(poolMonitor)
	}

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	return &Client{
		MongoClient: client,
		Database:    client.Database(database),
	}, nil
}

func (c *Client) Collection(name string) *mongo.Collection {
	return c.Database.Collection(name)
}

func (c *Client) Disconnect(ctx context.Context) error {
	return c.MongoClient.Disconnect(ctx)
}

func (c *Client) Ping(ctx context.Context) error {
	return c.MongoClient.Ping(ctx, nil)
}
