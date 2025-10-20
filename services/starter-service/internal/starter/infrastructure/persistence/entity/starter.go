package entity

import "time"

type StarterEntity struct {
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement"`
	Domain        string     `gorm:"column:domain;type:varchar(25);uniqueIndex;not null"`
	Name          string     `gorm:"column:name;type:varchar(255);not null"`
	Email         string     `gorm:"column:email;type:varchar(100)"`
	Mobile        string     `gorm:"column:mobile;type:varchar(20);not null"`
	WorkPhone     string     `gorm:"column:work_phone;type:varchar(20)"`
	JobTitle      string     `gorm:"column:job_title;type:varchar(100);not null"`
	DepartmentID  *int64     `gorm:"column:department_id"`
	LineManagerID *int64     `gorm:"column:line_manager_id"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at;index"`
}

func (StarterEntity) TableName() string {
	return "starters"
}
