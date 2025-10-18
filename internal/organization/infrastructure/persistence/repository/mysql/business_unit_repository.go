package mysql

import (
	"context"

	"github.com/kiin21/go-rest/internal/organization/domain/model"
	repo "github.com/kiin21/go-rest/internal/organization/domain/repository"
	"github.com/kiin21/go-rest/internal/organization/infrastructure/persistence/entity"
	sharedDomain "github.com/kiin21/go-rest/internal/shared/domain/model"
	"github.com/kiin21/go-rest/pkg/response"
	"gorm.io/gorm"
)

type MySQLBusinessUnitRepository struct {
	db *gorm.DB
}

func NewMySQLBusinessUnitRepository(db *gorm.DB) repo.BusinessUnitRepository {
	return &MySQLBusinessUnitRepository{db: db}
}

func (r *MySQLBusinessUnitRepository) FindByID(ctx context.Context, id int64) (*model.BusinessUnit, error) {
	var model entity.BusinessUnitModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&model), nil
}

func (r *MySQLBusinessUnitRepository) FindByIDs(ctx context.Context, ids []int64) ([]*model.BusinessUnit, error) {
	if len(ids) == 0 {
		return []*model.BusinessUnit{}, nil
	}

	var models []entity.BusinessUnitModel
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&models).Error; err != nil {
		return nil, err
	}

	domains := make([]*model.BusinessUnit, len(models))
	for i, m := range models {
		domains[i] = r.toDomain(&m)
	}
	return domains, nil
}

func (r *MySQLBusinessUnitRepository) List(ctx context.Context, pg response.ReqPagination) ([]*model.BusinessUnit, int64, error) {
	var models []entity.BusinessUnitModel
	var total *int64

	query := r.db.WithContext(ctx).Model(&entity.BusinessUnitModel{})

	if err := query.Count(total).Error; err != nil {
		return nil, 0, err
	}

	offset := (pg.Page - 1) * pg.Limit
	if err := query.Offset(offset).Limit(pg.Limit).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	units := make([]*model.BusinessUnit, len(models))
	for i, m := range models {
		units[i] = r.toDomain(&m)
	}

	return units, *total, nil
}

func (r *MySQLBusinessUnitRepository) FindByIDWithDetails(ctx context.Context, id int64) (*model.BusinessUnitWithDetails, error) {
	var _model entity.BusinessUnitModel
	if err := r.db.WithContext(ctx).
		Preload("Company").
		Preload("Leader").
		First(&_model, id).
		Error; err != nil {
		return nil, err
	}
	return r.toDomainWithDetails(&_model), nil
}

func (r *MySQLBusinessUnitRepository) ListWithDetails(ctx context.Context, pg response.ReqPagination) ([]*model.BusinessUnitWithDetails, int64, error) {
	var models []entity.BusinessUnitModel
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.BusinessUnitModel{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (pg.Page - 1) * pg.Limit
	if offset < 0 {
		offset = 0
	}

	if err := query.
		Preload("Company").
		Preload("Leader").
		Order("name ASC").
		Offset(offset).Limit(pg.Limit).
		Find(&models).
		Error; err != nil {
		return nil, 0, err
	}

	units := make([]*model.BusinessUnitWithDetails, len(models))
	for i, m := range models {
		units[i] = r.toDomainWithDetails(&m)
	}

	return units, total, nil
}

// =============== UTILS ===================
func (r *MySQLBusinessUnitRepository) toDomainWithDetails(m *entity.BusinessUnitModel) *model.BusinessUnitWithDetails {
	bu := &model.BusinessUnitWithDetails{
		BusinessUnit: r.toDomain(m), // Reuse existing converter
	}

	if m.Company != nil {
		bu.Company = &model.Company{
			ID:   m.Company.ID,
			Name: m.Company.Name,
		}
	}

	if m.Leader != nil {
		bu.Leader = &sharedDomain.LineManagerNested{
			ID:       m.Leader.ID,
			Domain:   m.Leader.Domain,
			Name:     m.Leader.Name,
			Email:    m.Leader.Email,
			JobTitle: m.Leader.JobTitle,
		}
	}

	return bu
}

func (r *MySQLBusinessUnitRepository) toDomain(m *entity.BusinessUnitModel) *model.BusinessUnit {
	return &model.BusinessUnit{
		ID:        m.ID,
		Name:      m.Name,
		Shortname: m.Shortname,
		CompanyID: m.CompanyID,
		LeaderID:  m.LeaderID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
