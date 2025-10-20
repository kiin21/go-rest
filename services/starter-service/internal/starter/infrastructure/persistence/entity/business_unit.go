package entity

import (
	"time"

	"gorm.io/gorm"
)

type BusinessUnitEntity struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string         `gorm:"column:name;not null"`
	Shortname string         `gorm:"column:shortname;not null"`
	CompanyID int64          `gorm:"column:company_id;not null"`
	LeaderID  *int64         `gorm:"column:leader_id"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`

	// Relationships
	Company *CompanyEntity `gorm:"foreignKey:CompanyID"`
	Leader  *StarterEntity `gorm:"foreignKey:LeaderID"`
}

func (BusinessUnitEntity) TableName() string {
	return "business_units"
}
