package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/kiin21/go-rest/pkg/httputil"
	domainModel "github.com/kiin21/go-rest/services/notification-service/internal/notification/domain/model"
	domainRepo "github.com/kiin21/go-rest/services/notification-service/internal/notification/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultSortBy    = "timestamp"
	defaultSortOrder = "desc"
)

type notificationMongoRepository struct {
	collection *mongo.Collection
}

func NewNotificationMongoRepository(collection *mongo.Collection) domainRepo.NotificationRepository {
	return &notificationMongoRepository{
		collection: collection,
	}
}

func (r *notificationMongoRepository) List(
	ctx context.Context,
	filter domainRepo.ListNotificationsFilter,
	pagination httputil.ReqPagination,
) ([]*domainModel.Notification, int64, error) {
	if pagination.Page <= 0 {
		pagination.Page = 1
	}
	if pagination.Limit <= 0 {
		pagination.Limit = 20
	}

	sortBy := mapSortField(filter.SortBy)

	sortOrder := 1
	if filter.SortOrder != "" {
		if strings.EqualFold(filter.SortOrder, "desc") {
			sortOrder = -1
		}
	} else if defaultSortOrder == "desc" {
		sortOrder = -1
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: sortBy, Value: sortOrder}})
	findOptions.SetSkip(int64((pagination.Page - 1) * pagination.Limit))
	findOptions.SetLimit(int64(pagination.Limit))

	cursor, err := r.collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	results := make([]*domainModel.Notification, 0, pagination.Limit)
	for cursor.Next(ctx) {
		var doc notificationDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, 0, err
		}
		results = append(results, doc.toDomain())
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (r *notificationMongoRepository) Create(ctx context.Context, notification *domainModel.Notification) error {
	if notification == nil {
		return errors.New("notification is nil")
	}

	ts := notification.Timestamp
	if ts.IsZero() {
		ts = time.Now().UTC()
	}

	doc := notificationDocument{
		ID:          primitive.NewObjectID(),
		FromStarter: notification.FromStarter,
		ToStarter:   notification.ToStarter,
		Message:     notification.Message,
		Type:        notification.Type,
		Timestamp:   ts,
	}

	if _, err := r.collection.InsertOne(ctx, doc); err != nil {
		return err
	}

	notification.ID = doc.ID.Hex()
	notification.Timestamp = ts
	return nil
}

type notificationDocument struct {
	ID          primitive.ObjectID `bson:"_id"`
	FromStarter string             `bson:"from_starter"`
	ToStarter   string             `bson:"to_starter"`
	Message     string             `bson:"message"`
	Type        string             `bson:"type"`
	Timestamp   time.Time          `bson:"timestamp"`
}

func (d *notificationDocument) toDomain() *domainModel.Notification {
	return &domainModel.Notification{
		ID:          d.ID.Hex(),
		FromStarter: d.FromStarter,
		ToStarter:   d.ToStarter,
		Message:     d.Message,
		Type:        d.Type,
		Timestamp:   d.Timestamp,
	}
}

func mapSortField(input string) string {
	switch strings.ToLower(input) {
	case "from":
		return "from_starter"
	case "to":
		return "to_starter"
	case "timestamp":
		return "timestamp"
	case "type":
		return "type"
	default:
		return defaultSortBy
	}
}
