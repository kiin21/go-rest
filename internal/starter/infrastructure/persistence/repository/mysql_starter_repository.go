package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	sharedDomain "github.com/kiin21/go-rest/internal/shared/domain"
	"github.com/kiin21/go-rest/internal/shared/infrastructure/persistence/model"
	starterAggregate "github.com/kiin21/go-rest/internal/starter/domain/aggregate"
	starterPort "github.com/kiin21/go-rest/internal/starter/domain/port"
	"github.com/kiin21/go-rest/pkg/response"
	"gorm.io/gorm"
)

// MySQLStarterRepository implements StarterRepository using MySQL
type MySQLStarterRepository struct {
	db *gorm.DB
}

// NewMySQLStarterRepository creates a new MySQL repository
func NewMySQLStarterRepository(db *gorm.DB) starterPort.StarterRepository {

	return &MySQLStarterRepository{db: db}
}

// Create creates a new starter in the database
// Note: Domain uniqueness should be validated by Domain Service before calling this
func (r *MySQLStarterRepository) Create(ctx context.Context, starter *starterAggregate.Starter) error {
	model := r.toModel(starter)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return sharedDomain.ErrDuplicateEntry
		}
		return fmt.Errorf("failed to create starter: %w", err)
	}

	return nil
}

// Update updates an existing starter
func (r *MySQLStarterRepository) Update(ctx context.Context, starter *starterAggregate.Starter) error {
	starterModel := r.toModel(starter)

	result := r.db.WithContext(ctx).
		Model(&model.StarterModel{}).
		Where("id = ?", starterModel.ID).
		Updates(starterModel)

	if result.Error != nil {
		return fmt.Errorf("failed to update starter: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return sharedDomain.ErrNotFound
	}

	return nil
}

func (r *MySQLStarterRepository) FindByID(ctx context.Context, id int64) (*starterAggregate.Starter, error) {
	var model model.StarterModel

	if err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedDomain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find starter by id: %w", err)
	}

	return r.toDomain(&model), nil
}

func (r *MySQLStarterRepository) FindByDomain(ctx context.Context, domainName string) (*starterAggregate.Starter, error) {
	var model model.StarterModel

	if err := r.db.WithContext(ctx).Where("domain = ? AND deleted_at IS NULL", domainName).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedDomain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find starter by domain: %w", err)
	}

	return r.toDomain(&model), nil
}

// SoftDelete - soft deletes a starter by setting deleted_at timestamp
func (r *MySQLStarterRepository) SoftDelete(ctx context.Context, domainName string) error {
	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&model.StarterModel{}).
		Where("domain = ? AND deleted_at IS NULL", domainName).
		Update("deleted_at", now)

	if result.Error != nil {
		return fmt.Errorf("failed to soft delete starter: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return sharedDomain.ErrNotFound
	}

	return nil
}

// List returns a paginated list of starters with optional filters (excludes soft deleted)
func (r *MySQLStarterRepository) List(ctx context.Context, filter starterPort.ListFilter, pg response.ReqPagination) ([]*starterAggregate.Starter, int64, error) {
	var models []model.StarterModel
	var total int64

	query := r.db.WithContext(ctx).Model(&model.StarterModel{}).Where("deleted_at IS NULL")

	// Apply filters
	query = r.applyFilters(query, filter)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count starters: %w", err)
	}

	// Apply pagination
	offset := (pg.Page - 1) * pg.Limit
	if err := query.Offset(offset).Limit(pg.Limit).Find(&models).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list starters: %w", err)
	}

	starters := make([]*starterAggregate.Starter, len(models))
	for i, model := range models {
		starters[i] = r.toDomain(&model)
	}

	return starters, total, nil
}

// SearchByKeyword performs keyword search across domain, email, mobile, job title (excludes soft deleted)
func (r *MySQLStarterRepository) SearchByKeyword(ctx context.Context, keyword string, filter starterPort.ListFilter, pg response.ReqPagination) ([]*starterAggregate.Starter, int64, error) {
	var models []model.StarterModel
	var total int64

	query := r.db.WithContext(ctx).Model(&model.StarterModel{}).Where("deleted_at IS NULL")

	// Apply keyword search
	searchPattern := "%" + keyword + "%"
	query = query.Where(
		"domain LIKE ? OR email LIKE ? OR mobile LIKE ? OR job_title LIKE ?",
		searchPattern, searchPattern, searchPattern, searchPattern,
	)

	// Apply additional filters
	query = r.applyFilters(query, filter)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// Apply pagination
	offset := (pg.Page - 1) * pg.Limit
	if err := query.Offset(offset).Limit(pg.Limit).Find(&models).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search starters: %w", err)
	}

	starters := make([]*starterAggregate.Starter, len(models))
	for i, model := range models {
		starters[i] = r.toDomain(&model)
	}

	return starters, total, nil
}

// FindSubordinates finds all subordinates of a line manager
func (r *MySQLStarterRepository) FindSubordinates(ctx context.Context, managerID int64) ([]*starterAggregate.Starter, error) {
	var models []model.StarterModel

	if err := r.db.WithContext(ctx).Where("line_manager_id = ?", managerID).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find subordinates: %w", err)
	}

	starters := make([]*starterAggregate.Starter, len(models))
	for i, model := range models {
		starters[i] = r.toDomain(&model)
	}

	return starters, nil
}

// applyFilters applies ListFilter to query
func (r *MySQLStarterRepository) applyFilters(query *gorm.DB, filter starterPort.ListFilter) *gorm.DB {
	if filter.DepartmentID != nil {
		query = query.Where("department_id = ?", *filter.DepartmentID)
	}

	if filter.BusinessUnitID != nil {
		// This would typically join with department/business_unit table
		// Simplified: assuming we have business_unit_id column
		// query = query.Where("business_unit_id = ?", *filter.BusinessUnitID)
	}

	if filter.LineManagerID != nil {
		query = query.Where("line_manager_id = ?", *filter.LineManagerID)
	}

	return query
}

// toModel converts domain Starter to database entity
func (r *MySQLStarterRepository) toModel(starter *starterAggregate.Starter) *model.StarterModel {
	return &model.StarterModel{
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

// toDomain converts database entity to domain Starter
func (r *MySQLStarterRepository) toDomain(model *model.StarterModel) *starterAggregate.Starter {
	return starterAggregate.Rehydrate(
		model.ID,
		model.Domain,
		model.Name,
		model.Email,
		model.Mobile,
		model.WorkPhone,
		model.JobTitle,
		model.DepartmentID,
		model.LineManagerID,
		model.CreatedAt,
		model.UpdatedAt,
	)
}
