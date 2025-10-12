package repository

import (
	"context"

	domain "github.com/kiin21/go-rest/internal/organization/domain"
	model "github.com/kiin21/go-rest/internal/organization/infrastructure/persistence/model"
	"gorm.io/gorm"
)

// MySQLBusinessUnitRepository implements BusinessUnitRepository
type MySQLBusinessUnitRepository struct {
	db *gorm.DB
}

// NewMySQLBusinessUnitRepository creates repository
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

	units := make([]*domain.BusinessUnit, len(models))
	for i, m := range models {
		units[i] = r.toDomain(&m)
	}
	return units, nil
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
