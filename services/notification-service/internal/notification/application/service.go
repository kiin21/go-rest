package application

import (
	"context"
	"math"
	"strconv"

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

	value := strconv.Itoa(pagination.GetPage() - 1)
	prev = &value

	totalPages := 0
	totalPages = int(math.Ceil(float64(total) / float64(pagination.GetLimit())))

	if totalPages > 0 && pagination.GetPage() < totalPages {
		value := strconv.Itoa(pagination.GetPage() + 1)
		next = &value
	}

	return &httputil.PaginatedResult[*domainmodel.Notification]{
		Data: data,
		Pagination: httputil.RespPagination{
			Limit:      pagination.GetLimit(),
			TotalItems: total,
			Prev:       prev,
			Next:       next,
		},
	}, nil
}
