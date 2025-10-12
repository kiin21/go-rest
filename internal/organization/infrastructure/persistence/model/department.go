package model

import "time"

// DepartmentModel is the GORM model for the departments table
type DepartmentModel struct {
	ID                int64      `gorm:"column:id;primaryKey"`
	GroupDepartmentID *int64     `gorm:"column:group_department_id;index"`
	FullName          string     `gorm:"column:full_name"`
	Shortname         string     `gorm:"column:shortname"`
	BusinessUnitID    *int64     `gorm:"column:business_unit_id;index:idx_departments_bu_deleted"`
	LeaderID          *int64     `gorm:"column:leader_id"`
	CreatedAt         time.Time  `gorm:"column:created_at"`
	UpdatedAt         time.Time  `gorm:"column:updated_at"`
	DeletedAt         *time.Time `gorm:"column:deleted_at;index:idx_departments_deleted_at,idx_departments_bu_deleted"`
}

// TableName specifies the table name for GORM
func (DepartmentModel) TableName() string {
	return "departments"
}