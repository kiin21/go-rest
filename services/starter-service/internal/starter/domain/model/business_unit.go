package model

import "time"

type BusinessUnit struct {
	ID        int64
	Name      string
	Shortname string
	CompanyID int64
	LeaderID  *int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
