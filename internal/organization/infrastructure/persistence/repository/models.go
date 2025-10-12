package repository

import (
	"time"

	"gorm.io/gorm"
)

// DepartmentModel represents the departments table in the database
type DepartmentModel struct {
	ID               int64          `gorm:"column:id;primaryKey;autoIncrement"`
	FullName         string         `gorm:"column:full_name;not null"`
	Shortname        string         `gorm:"column:shortname;not null"`
	GroupDepartment  *int64         `gorm:"column:group_department"`
	BusinessUnitID   *int64         `gorm:"column:business_unit_id"`
	DepartmentStatus int            `gorm:"column:department_status;not null;default:1"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at;index"`
	CreatedAt        time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        time.Time      `gorm:"column:updated_at;autoUpdateTime"`
}

func (DepartmentModel) TableName() string {
	return "departments"
}

// BusinessUnitModel represents the business_units table in the database
type BusinessUnitModel struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string         `gorm:"column:name;not null"`
	Shortname string         `gorm:"column:shortname;not null"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
}

func (BusinessUnitModel) TableName() string {
	return "business_units"
}
