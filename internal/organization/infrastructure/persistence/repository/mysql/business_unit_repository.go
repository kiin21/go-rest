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

type BusinessUnitRepository struct {
	db *gorm.DB
}

func NewBusinessUnitRepository(db *gorm.DB) repo.BusinessUnitRepository {
	return &BusinessUnitRepository{db: db}
}

func (r *BusinessUnitRepository) FindByIDs(ctx context.Context, ids []int64) ([]*model.BusinessUnit, error) {
	if len(ids) == 0 {
		return []*model.BusinessUnit{}, nil
	}

	var entities []entity.BusinessUnitEntity
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&entities).Error; err != nil {
		return nil, err
	}

	domains := make([]*model.BusinessUnit, len(entities))
	for i, m := range entities {
		domains[i] = r.toModel(&m)
	}
	return domains, nil
}

func (r *BusinessUnitRepository) List(ctx context.Context, pg response.ReqPagination) ([]*model.BusinessUnit, int64, error) {
	var entities []entity.BusinessUnitEntity
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.BusinessUnitEntity{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (pg.Page - 1) * pg.Limit
	if err := query.Offset(offset).Limit(pg.Limit).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	units := make([]*model.BusinessUnit, len(entities))
	for i, m := range entities {
		units[i] = r.toModel(&m)
	}

	return units, total, nil
}

func (r *BusinessUnitRepository) FindByIDWithDetails(ctx context.Context, id int64) (*model.BusinessUnitWithDetails, error) {
	var businessUnitEntity entity.BusinessUnitEntity
	if err := r.db.WithContext(ctx).
		Preload("Company").
		Preload("Leader").
		First(&businessUnitEntity, id).
		Error; err != nil {
		return nil, err
	}
	return r.toModelWithDetails(&businessUnitEntity), nil
}

func (r *BusinessUnitRepository) ListWithDetails(ctx context.Context, pg response.ReqPagination) ([]*model.BusinessUnitWithDetails, int64, error) {
	var models []entity.BusinessUnitEntity
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.BusinessUnitEntity{})

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
		units[i] = r.toModelWithDetails(&m)
	}

	return units, total, nil
}

// =============== UTILS ===================
func (r *BusinessUnitRepository) toModelWithDetails(m *entity.BusinessUnitEntity) *model.BusinessUnitWithDetails {
	bu := &model.BusinessUnitWithDetails{
		BusinessUnit: r.toModel(m), // Reuse existing converter
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

func (r *BusinessUnitRepository) toModel(m *entity.BusinessUnitEntity) *model.BusinessUnit {
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
