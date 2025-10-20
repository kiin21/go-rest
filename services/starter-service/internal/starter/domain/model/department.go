package model

import "time"

type Department struct {
	ID                int64
	GroupDepartmentID *int64
	FullName          string
	Shortname         string
	BusinessUnitID    *int64
	LeaderID          *int64
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
}

type DepartmentListFilter struct {
	BusinessUnitID *int64
}
