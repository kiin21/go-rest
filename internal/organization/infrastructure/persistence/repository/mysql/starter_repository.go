package mysql

import (
	"context"
	"errors"
	"time"

	model "github.com/kiin21/go-rest/internal/organization/domain/model"
	repo "github.com/kiin21/go-rest/internal/organization/domain/repository"
	sharedDomain "github.com/kiin21/go-rest/internal/shared/domain"
	"github.com/kiin21/go-rest/pkg/response"
	"gorm.io/gorm"
)

type MySQLStarterRepository struct {
	db *gorm.DB
}

func NewMySQLStarterRepository(db *gorm.DB) repo.StarterRepository {
	return &MySQLStarterRepository{db: db}
}

type StarterModel struct {
	ID            int64      `gorm:"primaryKey;column:id"`
	Domain        string     `gorm:"column:domain;uniqueIndex;not null"`
	Name          string     `gorm:"column:name;not null"`
	Email         string     `gorm:"column:email"`
	Mobile        string     `gorm:"column:mobile;not null"`
	WorkPhone     string     `gorm:"column:work_phone"`
	JobTitle      string     `gorm:"column:job_title;not null"`
	DepartmentID  *int64     `gorm:"column:department_id"`
	LineManagerID *int64     `gorm:"column:line_manager_id"`
	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt     *time.Time `gorm:"column:deleted_at;index"`
}

func (StarterModel) TableName() string {
	return "starters"
}

func (r *MySQLStarterRepository) FindByID(ctx context.Context, id int64) (*model.Starter, error) {
	var model StarterModel
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedDomain.ErrNotFound
		}
		return nil, err
	}

	return r.toDomainAggregate(&model), nil
}

func (r *MySQLStarterRepository) FindByDomain(ctx context.Context, domain string) (*model.Starter, error) {
	var model StarterModel
	err := r.db.WithContext(ctx).Where("domain = ? AND deleted_at IS NULL", domain).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedDomain.ErrNotFound
		}
		return nil, err
	}

	return r.toDomainAggregate(&model), nil
}

func (r *MySQLStarterRepository) List(ctx context.Context, filter model.StarterListFilter, pg response.ReqPagination) ([]*model.Starter, int64, error) {
	query := r.db.WithContext(ctx).Model(&StarterModel{}).Where("deleted_at IS NULL")

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

	var models []StarterModel
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	starters := make([]*model.Starter, len(models))
	for i, model := range models {
		starters[i] = r.toDomainAggregate(&model)
	}

	return starters, total, nil
}

func (r *MySQLStarterRepository) SearchByKeyword(ctx context.Context, keyword string, filter model.StarterListFilter, pg response.ReqPagination) ([]*model.Starter, int64, error) {
	query := r.db.WithContext(ctx).Model(&StarterModel{}).Where("deleted_at IS NULL")

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

	var models []StarterModel
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	starters := make([]*model.Starter, len(models))
	for i, model := range models {
		starters[i] = r.toDomainAggregate(&model)
	}

	return starters, total, nil
}

func (r *MySQLStarterRepository) Create(ctx context.Context, starter *model.Starter) error {
	model := r.fromDomainAggregate(starter)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *MySQLStarterRepository) Update(ctx context.Context, starter *model.Starter) error {
	model := r.fromDomainAggregate(starter)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *MySQLStarterRepository) SoftDelete(ctx context.Context, domain string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&StarterModel{}).
		Where("domain = ? AND deleted_at IS NULL", domain).
		Update("deleted_at", now).Error
}

// Helper methods for domain conversion
func (r *MySQLStarterRepository) toDomainAggregate(sm *StarterModel) *model.Starter {
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

func (r *MySQLStarterRepository) fromDomainAggregate(starter *model.Starter) *StarterModel {
	return &StarterModel{
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
