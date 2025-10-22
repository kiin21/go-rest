package mysql

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	starterquery "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/starter/query"
	sharedDomain "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/error"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
	repo "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/infrastructure/persistence/entity"
	"gorm.io/gorm"
)

type StarterRepository struct {
	db *gorm.DB
}

func NewStarterRepository(db *gorm.DB) repo.StarterRepository {
	return &StarterRepository{db: db}
}

func (r *StarterRepository) FindByID(ctx context.Context, id int64) (*model.Starter, error) {
	var starterEntity entity.StarterEntity
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&starterEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedDomain.ErrNotFound
		}
		return nil, err
	}

	return r.toModel(&starterEntity)
}

func (r *StarterRepository) FindByDomain(ctx context.Context, domain string) (*model.Starter, error) {
	var starterEntity entity.StarterEntity
	err := r.db.WithContext(ctx).Where("domain = ? AND deleted_at IS NULL", domain).First(&starterEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedDomain.ErrNotFound
		}
		return nil, err
	}

	return r.toModel(&starterEntity)
}

func (r *StarterRepository) SearchByKeyword(ctx context.Context, listStarterQuery starterquery.ListStartersQuery) ([]*model.Starter, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.StarterEntity{}).Where("deleted_at IS NULL")

	// Apply keyword search
	if listStarterQuery.Keyword != "" && listStarterQuery.SearchBy != "" {
		query = query.Where("? LIKE %?%", listStarterQuery.SearchBy, listStarterQuery.Keyword)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = r.applySort(query, listStarterQuery.SortBy, listStarterQuery.SortOrder)

	// Apply pagination
	offset := (listStarterQuery.Pagination.Page - 1) * listStarterQuery.Pagination.Limit
	query = query.Offset(offset).Limit(listStarterQuery.Pagination.Limit)

	var models []entity.StarterEntity
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	starters := make([]*model.Starter, 0, len(models))
	for i := range models {
		starterModel, err := r.toModel(&models[i])
		if err != nil {
			return nil, 0, err
		}
		starters = append(starters, starterModel)
	}

	return starters, total, nil
}

func (r *StarterRepository) Create(ctx context.Context, starter *model.Starter) error {
	domainAggregate := r.toEntity(starter)
	return r.db.WithContext(ctx).Create(domainAggregate).Error
}

func (r *StarterRepository) Update(ctx context.Context, starter *model.Starter) error {
	domainAggregate := r.toEntity(starter)
	return r.db.WithContext(ctx).Save(domainAggregate).Error
}

func (r *StarterRepository) SoftDelete(ctx context.Context, domain string) (*model.Starter, error) {
	var starterEntity entity.StarterEntity

	now := time.Now()
	err := r.db.WithContext(ctx).Model(&entity.StarterEntity{}).
		Where("domain = ? AND deleted_at IS NULL", domain).
		Update("deleted_at", now).First(&starterEntity).Error

	if err != nil {
		return nil, err
	}

	return r.toModel(&starterEntity)
}

// Helper methods for domain conversion
func (r *StarterRepository) toModel(sm *entity.StarterEntity) (*model.Starter, error) {
	return model.Rehydrate(
		sm.ID,
		sm.Domain,
		sm.Name,
		sm.Email,
		sm.Mobile,
		sm.WorkPhone,
		sm.JobTitle,
		sm.DepartmentID,
		sm.LineManagerID,
		sm.CreatedAt,
		sm.UpdatedAt,
	)
}

func (r *StarterRepository) toEntity(starter *model.Starter) *entity.StarterEntity {
	return &entity.StarterEntity{
		ID:            starter.ID(),
		Domain:        starter.Domain(),
		Name:          starter.Name(),
		Email:         starter.Email(),
		Mobile:        starter.Mobile(),
		WorkPhone:     starter.WorkPhone(),
		JobTitle:      starter.JobTitle(),
		DepartmentID:  starter.DepartmentID(),
		LineManagerID: starter.LineManagerID(),
		CreatedAt:     starter.CreatedAt(),
		UpdatedAt:     starter.UpdatedAt(),
	}
}

func (r *StarterRepository) applySort(query *gorm.DB, sortBy, sortOrder string) *gorm.DB {
	column := mapStarterSortColumn(sortBy)
	direction := strings.ToLower(sortOrder)
	if direction != "desc" {
		direction = "asc"
	}

	return query.Order(fmt.Sprintf("%s %s", column, direction))
}

func mapStarterSortColumn(sortBy string) string {
	switch strings.ToLower(sortBy) {
	case "domain":
		return "starters.domain"
	case "created_at":
		return "starters.created_at"
	default:
		return "starters.id"
	}
}
