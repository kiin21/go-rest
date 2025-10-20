package application

import (
	"context"
	"math"
	"strconv"
	"time"

	"github.com/kiin21/go-rest/pkg/events"
	"github.com/kiin21/go-rest/pkg/httputil"
	domainmodel "github.com/kiin21/go-rest/services/notification-service/internal/notification/domain/model"
	domainrepo "github.com/kiin21/go-rest/services/notification-service/internal/notification/domain/repository"
)

type NotiApplicationService struct {
	repo domainrepo.NotificationRepository
}

func NewNotiApplicationService(repo domainrepo.NotificationRepository) *NotiApplicationService {
	return &NotiApplicationService{repo: repo}
}

func (s *NotiApplicationService) ListNotifications(
	ctx context.Context,
	query ListNotificationsQuery,
) (*httputil.PaginatedResult[*domainmodel.Notification], error) {
	pagination := query.Pagination
	if pagination.Page <= 0 {
		pagination.Page = 1
	}
	if pagination.Limit <= 0 {
		pagination.Limit = 20
	}

	filter := domainrepo.ListNotificationsFilter{
		SortBy:    query.SortBy,
		SortOrder: query.SortOrder,
	}

	data, total, err := s.repo.List(ctx, filter, pagination)
	if err != nil {
		return nil, err
	}

	var (
		prev *string
		next *string
	)

	if pagination.Page > 1 {
		value := strconv.Itoa(pagination.Page - 1)
		prev = &value
	}

	totalPages := 0
	if pagination.Limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(pagination.Limit)))
	}

	if totalPages > 0 && pagination.Page < totalPages {
		value := strconv.Itoa(pagination.Page + 1)
		next = &value
	}

	return &httputil.PaginatedResult[*domainmodel.Notification]{
		Data: data,
		Pagination: httputil.RespPagination{
			Limit:      pagination.Limit,
			TotalItems: total,
			Prev:       prev,
			Next:       next,
		},
	}, nil
}

func (s *NotiApplicationService) StoreNotification(ctx context.Context, event *events.LeaderAssignmentNotification) error {
	if event == nil {
		return nil
	}

	timestamp := event.Timestamp
	if timestamp.IsZero() {
		timestamp = time.Now().UTC()
	}

	notification := &domainmodel.Notification{
		FromStarter: event.FromStarter,
		ToStarter:   event.ToStarter,
		Message:     event.Message,
		Type:        event.Type,
		Timestamp:   timestamp,
	}

	return s.repo.Create(ctx, notification)
}
