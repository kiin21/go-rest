package model

import "time"

// BusinessUnitModel is the GORM model for business_units table
type BusinessUnitModel struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name"`
	Shortname string    `gorm:"column:shortname"`
	CompanyID int64     `gorm:"column:company_id"`
	LeaderID  *int64    `gorm:"column:leader_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// TableName specifies table name
func (BusinessUnitModel) TableName() string {
	return "business_units"
}
