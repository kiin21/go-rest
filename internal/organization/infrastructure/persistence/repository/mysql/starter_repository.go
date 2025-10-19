package mysql

import (
	"context"
	"errors"
	"time"

	"github.com/kiin21/go-rest/internal/organization/domain/model"
	repo "github.com/kiin21/go-rest/internal/organization/domain/repository"
	"github.com/kiin21/go-rest/internal/organization/infrastructure/persistence/entity"
	sharedDomain "github.com/kiin21/go-rest/internal/shared/domain"
	"github.com/kiin21/go-rest/pkg/response"
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

	return r.toModel(&starterEntity), nil
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

	return r.toModel(&starterEntity), nil
}

func (r *StarterRepository) List(ctx context.Context, filter model.StarterListFilter, pg response.ReqPagination) ([]*model.Starter, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.StarterEntity{}).Where("deleted_at IS NULL")

	// Apply filters
	if filter.DepartmentID != nil {
		query = query.Where("department_id = ?", *filter.DepartmentID)
	}
	if filter.BusinessUnitID != nil {
		query = query.Joins("JOIN departments ON starters.department_id = departments.id").
			Where("departments.business_unit_id = ?", *filter.BusinessUnitID)
	}
	if filter.LineManagerID != nil {
		query = query.Where("line_manager_id = ?", *filter.LineManagerID)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (pg.Page - 1) * pg.Limit
	query = query.Offset(offset).Limit(pg.Limit)

	var models []entity.StarterEntity
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	starters := make([]*model.Starter, len(models))
	for i, starterEntity := range models {
		starters[i] = r.toModel(&starterEntity)
	}

	return starters, total, nil
}

func (r *StarterRepository) SearchByKeyword(ctx context.Context, keyword string, filter model.StarterListFilter, pg response.ReqPagination) ([]*model.Starter, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.StarterEntity{}).Where("deleted_at IS NULL")

	// Apply keyword search
	searchPattern := "%" + keyword + "%"
	query = query.Where("domain LIKE ? OR name LIKE ? OR email LIKE ?", searchPattern, searchPattern, searchPattern)

	// Apply filters
	if filter.DepartmentID != nil {
		query = query.Where("department_id = ?", *filter.DepartmentID)
	}
	if filter.BusinessUnitID != nil {
		query = query.Joins("JOIN departments ON starters.department_id = departments.id").
			Where("departments.business_unit_id = ?", *filter.BusinessUnitID)
	}
	if filter.LineManagerID != nil {
		query = query.Where("line_manager_id = ?", *filter.LineManagerID)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (pg.Page - 1) * pg.Limit
	query = query.Offset(offset).Limit(pg.Limit)

	var models []entity.StarterEntity
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	starters := make([]*model.Starter, len(models))
	for i, starterEntity := range models {
		starters[i] = r.toModel(&starterEntity)
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

func (r *StarterRepository) SoftDelete(ctx context.Context, domain string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entity.StarterEntity{}).
		Where("domain = ? AND deleted_at IS NULL", domain).
		Update("deleted_at", now).Error
}

// Helper methods for domain conversion
func (r *StarterRepository) toModel(sm *entity.StarterEntity) *model.Starter {
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
