package repository

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/kiin21/go-rest/pkg/httputil"
	domainModel "github.com/kiin21/go-rest/services/notification-service/internal/notification/domain/model"
	domainRepo "github.com/kiin21/go-rest/services/notification-service/internal/notification/domain/repository"
	"github.com/kiin21/go-rest/services/notification-service/internal/notification/infrastructure/repository/document"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultSortBy = "timestamp"
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
	pg httputil.ReqPagination,
) ([]*domainModel.Notification, int64, error) {

	sortBy := mapSortField(filter.SortBy)
	sortOrder := mapSortOrder(filter.SortOrder)

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: sortBy, Value: sortOrder}})
	findOptions.SetSkip(int64(pg.GetOffset()))
	findOptions.SetLimit(int64(pg.GetLimit()))

	cursor, err := r.collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, 0, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("Error closing MongoDB cursor: %v", err)
		}
	}(cursor, ctx)

	results := make([]*domainModel.Notification, 0, pg.GetLimit())
	for cursor.Next(ctx) {
		var doc document.NotificationDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, 0, err
		}
		results = append(results, doc.ToDomain())
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

	// Convert domain model to document with BSON tags
	doc := document.FromDomain(notification)

	if _, err := r.collection.InsertOne(ctx, doc); err != nil {
		return err
	}
	return nil
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

func mapSortOrder(input string) int {
	switch strings.ToLower(input) {
	case "asc":
		return 1
	case "desc":
		return -1
	default:
		return -1 // default to descending
	}
}
