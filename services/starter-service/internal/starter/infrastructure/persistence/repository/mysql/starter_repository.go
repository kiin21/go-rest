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

func (r *StarterRepository) FindByIDs(ctx context.Context, ids []int64) ([]*model.Starter, error) {
	var starterEntities []entity.StarterEntity
	err := r.db.WithContext(ctx).
		Where("id IN ? AND deleted_at IS NULL", ids).
		Find(&starterEntities).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find starters by ids: %w", err)
	}

	// Convert entities to models
	starters := make([]*model.Starter, 0, len(starterEntities))
	for _, starterEntity := range starterEntities {
		starter, err := r.toModel(&starterEntity)
		if err != nil {
			return nil, fmt.Errorf("failed to convert starterEntity to model: %w", err)
		}
		starters = append(starters, starter)
	}

	return starters, nil
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
func (r *StarterRepository) SearchByKeyword(ctx context.Context, listStarterQuery *starterquery.ListStartersQuery) ([]*model.Starter, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.StarterEntity{}).Where("starters.deleted_at IS NULL")

	// Apply keyword search
	if listStarterQuery.Keyword != "" && listStarterQuery.SearchBy != "" {
		searchPattern := "%" + listStarterQuery.Keyword + "%"

		switch listStarterQuery.SearchBy {
		case "domain":
			query = query.Where("starters.domain LIKE ?", searchPattern)
		case "fullname":
			query = query.Where("starters.name LIKE ?", searchPattern)
		case "dept_name":
			query = query.Joins("LEFT JOIN departments ON departments.id = starters.department_id AND departments.deleted_at IS NULL").
				Where("departments.full_name LIKE ? OR departments.shortname LIKE ?", searchPattern, searchPattern)
		case "bu_name":
			query = query.Joins("LEFT JOIN departments ON departments.id = starters.department_id AND departments.deleted_at IS NULL").
				Joins("LEFT JOIN business_units ON business_units.id = departments.business_unit_id").
				Where("business_units.name LIKE ? OR business_units.shortname LIKE ?", searchPattern, searchPattern)
		}
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = r.applySort(query, listStarterQuery.SortBy, listStarterQuery.SortOrder)

	// Apply pagination
	offset := listStarterQuery.Pagination.GetOffset()
	limit := listStarterQuery.Pagination.GetLimit()
	query = query.Offset(offset).Limit(limit)

	var models []entity.StarterEntity
	if err := query.Select("starters.*").Find(&models).Error; err != nil {
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
	starterEntity := r.toEntity(starter)
	
	if err := r.db.WithContext(ctx).Create(starterEntity).Error; err != nil {
		return err
	}

	starter.ID = starterEntity.ID
	starter.CreatedAt = starterEntity.CreatedAt
	starter.UpdatedAt = starterEntity.UpdatedAt

	return nil
}

func (r *StarterRepository) Update(ctx context.Context, starter *model.Starter) error {
	starterEntity := r.toEntity(starter)
	return r.db.WithContext(ctx).Save(starterEntity).Error
}

func (r *StarterRepository) SoftDelete(ctx context.Context, domain string) (*model.Starter, error) {
	var starterEntity entity.StarterEntity

	err := r.db.WithContext(ctx).
		Where("domain = ? AND deleted_at IS NULL", domain).
		First(&starterEntity).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedDomain.ErrNotFound
		}
		return nil, err
	}

	now := time.Now()
	err = r.db.WithContext(ctx).
		Model(&starterEntity).
		Update("deleted_at", now).Error

	if err != nil {
		return nil, err
	}

	starterEntity.DeletedAt = &now

	return r.toModel(&starterEntity)
}

func (r *StarterRepository) toModel(e *entity.StarterEntity) (*model.Starter, error) {
	return model.Rehydrate(e.ID, e.Domain, e.Name, e.Email, e.Mobile, e.WorkPhone, e.JobTitle, e.DepartmentID, e.LineManagerID, e.CreatedAt, e.UpdatedAt)
}

func (r *StarterRepository) toEntity(starter *model.Starter) *entity.StarterEntity {
	return &entity.StarterEntity{
		ID:            starter.ID,
		Domain:        starter.Domain,
		Name:          starter.Name,
		Email:         starter.Email.Value(),
		Mobile:        starter.Mobile,
		WorkPhone:     starter.WorkPhone,
		JobTitle:      starter.JobTitle,
		DepartmentID:  starter.DepartmentID,
		LineManagerID: starter.LineManagerID,
		CreatedAt:     starter.CreatedAt,
		UpdatedAt:     starter.UpdatedAt,
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
