package entity

import (
	"time"

	sharedModel "github.com/kiin21/go-rest/internal/shared/infrastructure/persistence/model"
	"gorm.io/gorm"
)

type BusinessUnitModel struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string         `gorm:"column:name;not null"`
	Shortname string         `gorm:"column:shortname;not null"`
	CompanyID int64          `gorm:"column:company_id;not null"`
	LeaderID  *int64         `gorm:"column:leader_id"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`

	// Relationships
	Company *CompanyModel             `gorm:"foreignKey:CompanyID"`
	Leader  *sharedModel.StarterModel `gorm:"foreignKey:LeaderID"`
}

func (BusinessUnitModel) TableName() string {
	return "business_units"
}
