package repository

import (
	"context"

	"github.com/kiin21/go-rest/pkg/httputil"
	"github.com/kiin21/go-rest/services/notification-service/internal/notification/domain/model"
)

type ListNotificationsFilter struct {
	SortBy    string
	SortOrder string
}

type NotificationRepository interface {
	List(ctx context.Context, filter ListNotificationsFilter, pagination httputil.ReqPagination) ([]*model.Notification, int64, error)
	Create(ctx context.Context, notification *model.Notification) error
}
