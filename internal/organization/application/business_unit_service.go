package application

import (
	"context"
	"strconv"

	appDto "github.com/kiin21/go-rest/internal/organization/application/dto"
	"github.com/kiin21/go-rest/internal/organization/domain"
	"github.com/kiin21/go-rest/pkg/response"
)

// BusinessUnitApplicationService handles business unit related use cases.
type BusinessUnitApplicationService struct {
	businessUnitRepo domain.BusinessUnitRepository
}

// NewBusinessUnitApplicationService creates a new business unit application service.
func NewBusinessUnitApplicationService(repo domain.BusinessUnitRepository) *BusinessUnitApplicationService {
	return &BusinessUnitApplicationService{
		businessUnitRepo: repo,
	}
}

// ListBusinessUnits returns a paginated list of business units.
func (s *BusinessUnitApplicationService) ListBusinessUnits(ctx context.Context, query appDto.ListBusinessUnitsQuery) (*response.PaginatedResult[*domain.BusinessUnit], error) {
	units, total, err := s.businessUnitRepo.List(ctx, query.Pagination)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / query.Pagination.Limit
	if int(total)%query.Pagination.Limit > 0 {
		totalPages++
	}

	var prev, next *string
	if query.Pagination.Page > 1 {
		value := strconv.Itoa(query.Pagination.Page - 1)
		prev = &value
	}
	if query.Pagination.Page < totalPages {
		value := strconv.Itoa(query.Pagination.Page + 1)
		next = &value
	}

	return &response.PaginatedResult[*domain.BusinessUnit]{
		Data: units,
		Pagination: response.RespPagination{
			Limit:      query.Pagination.Limit,
			TotalItems: total,
			Prev:       prev,
			Next:       next,
		},
	}, nil
}

// GetBusinessUnit retrieves a single business unit by ID.
func (s *BusinessUnitApplicationService) GetBusinessUnit(ctx context.Context, id int64) (*domain.BusinessUnit, error) {
	return s.businessUnitRepo.FindByID(ctx, id)
}

// ListBusinessUnitsWithDetails returns a paginated list of business units with company and leader details.
func (s *BusinessUnitApplicationService) ListBusinessUnitsWithDetails(ctx context.Context, query appDto.ListBusinessUnitsQuery) (*response.PaginatedResult[*domain.BusinessUnitWithDetails], error) {
	units, total, err := s.businessUnitRepo.ListWithDetails(ctx, query.Pagination)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / query.Pagination.Limit
	if int(total)%query.Pagination.Limit > 0 {
		totalPages++
	}

	var prev, next *string
	if query.Pagination.Page > 1 {
		value := strconv.Itoa(query.Pagination.Page - 1)
		prev = &value
	}
	if query.Pagination.Page < totalPages {
		value := strconv.Itoa(query.Pagination.Page + 1)
		next = &value
	}

	return &response.PaginatedResult[*domain.BusinessUnitWithDetails]{
		Data: units,
		Pagination: response.RespPagination{
			Limit:      query.Pagination.Limit,
			TotalItems: total,
			Prev:       prev,
			Next:       next,
		},
	}, nil
}

// GetBusinessUnitWithDetails retrieves a single business unit with company and leader details by ID.
func (s *BusinessUnitApplicationService) GetBusinessUnitWithDetails(ctx context.Context, id int64) (*domain.BusinessUnitWithDetails, error) {
	return s.businessUnitRepo.FindByIDWithDetails(ctx, id)
}
