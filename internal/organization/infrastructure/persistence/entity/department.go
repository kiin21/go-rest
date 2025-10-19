package entity

import (
	"time"
)

type DepartmentEntity struct {
	ID                int64      `gorm:"column:id;primaryKey;autoIncrement"`
	FullName          string     `gorm:"column:full_name;not null"`
	Shortname         string     `gorm:"column:shortname;not null"`
	GroupDepartmentID *int64     `gorm:"column:group_department_id"`
	BusinessUnitID    *int64     `gorm:"column:business_unit_id"`
	LeaderID          *int64     `gorm:"column:leader_id;not null;default:1"`
	DeletedAt         *time.Time `gorm:"column:deleted_at;index"`
	CreatedAt         time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (DepartmentEntity) TableName() string {
	return "departments"
}
