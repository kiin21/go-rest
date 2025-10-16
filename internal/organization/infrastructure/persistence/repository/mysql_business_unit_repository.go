package repository

import (
	"context"

	"github.com/kiin21/go-rest/internal/organization/domain"
	"github.com/kiin21/go-rest/internal/organization/infrastructure/persistence/model"
	"github.com/kiin21/go-rest/pkg/response"
	"gorm.io/gorm"
)

type MySQLBusinessUnitRepository struct {
	db *gorm.DB
}

func NewMySQLBusinessUnitRepository(db *gorm.DB) domain.BusinessUnitRepository {
	return &MySQLBusinessUnitRepository{db: db}
}

// FindByID retrieves business unit by ID
func (r *MySQLBusinessUnitRepository) FindByID(ctx context.Context, id int64) (*domain.BusinessUnit, error) {
	var model model.BusinessUnitModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&model), nil
}

// FindByIDs batch retrieves business units
func (r *MySQLBusinessUnitRepository) FindByIDs(ctx context.Context, ids []int64) ([]*domain.BusinessUnit, error) {
	if len(ids) == 0 {
		return []*domain.BusinessUnit{}, nil
	}

	var models []model.BusinessUnitModel
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&models).Error; err != nil {
		return nil, err
	}

	domains := make([]*domain.BusinessUnit, len(models))
	for i, m := range models {
		domains[i] = r.toDomain(&m)
	}
	return domains, nil
}

// List retrieves business units with pagination support
func (r *MySQLBusinessUnitRepository) List(ctx context.Context, pg response.ReqPagination) ([]*domain.BusinessUnit, int64, error) {
	var models []model.BusinessUnitModel
	var total int64

	query := r.db.WithContext(ctx).Model(&model.BusinessUnitModel{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (pg.Page - 1) * pg.Limit
	if offset < 0 {
		offset = 0
	}

	if err := query.Order("name ASC").Offset(offset).Limit(pg.Limit).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	units := make([]*domain.BusinessUnit, len(models))
	for i, m := range models {
		units[i] = r.toDomain(&m)
	}

	return units, total, nil
}

// toDomain converts model to domain
func (r *MySQLBusinessUnitRepository) toDomain(m *model.BusinessUnitModel) *domain.BusinessUnit {
	return &domain.BusinessUnit{
		ID:        m.ID,
		Name:      m.Name,
		Shortname: m.Shortname,
		CompanyID: m.CompanyID,
		LeaderID:  m.LeaderID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// FindByIDWithDetails retrieves a business unit with its company and leader details.
func (r *MySQLBusinessUnitRepository) FindByIDWithDetails(ctx context.Context, id int64) (*domain.BusinessUnitWithDetails, error) {
	var _model model.BusinessUnitModel
	if err := r.db.WithContext(ctx).Preload("Company").Preload("Leader").First(&_model, id).Error; err != nil {
		return nil, err
	}
	return r.toDomainWithDetails(&_model), nil
}

// ListWithDetails retrieves a paginated list of business units with their company and leader details.
func (r *MySQLBusinessUnitRepository) ListWithDetails(ctx context.Context, pg response.ReqPagination) ([]*domain.BusinessUnitWithDetails, int64, error) {
	var models []model.BusinessUnitModel
	var total int64

	query := r.db.WithContext(ctx).Model(&model.BusinessUnitModel{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (pg.Page - 1) * pg.Limit
	if offset < 0 {
		offset = 0
	}

	if err := query.Preload("Company").Preload("Leader").Order("name ASC").Offset(offset).Limit(pg.Limit).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	units := make([]*domain.BusinessUnitWithDetails, len(models))
	for i, m := range models {
		units[i] = r.toDomainWithDetails(&m)
	}

	return units, total, nil
}

// toDomainWithDetails converts model to domain with details.
func (r *MySQLBusinessUnitRepository) toDomainWithDetails(m *model.BusinessUnitModel) *domain.BusinessUnitWithDetails {
	bu := &domain.BusinessUnitWithDetails{
		BusinessUnit: r.toDomain(m), // Reuse existing converter
	}

	if m.Company != nil {
		bu.Company = &domain.Company{
			ID:   m.Company.ID,
			Name: m.Company.Name,
		}
	}

	if m.Leader != nil {
		bu.Leader = &domain.Leader{
			ID:       m.Leader.ID,
			Domain:   m.Leader.Domain,
			Name:     m.Leader.Name,
			Email:    m.Leader.Email,
			JobTitle: m.Leader.JobTitle,
		}
	}

	return bu
}
