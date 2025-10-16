package domain

import "time"

// BusinessUnit represents a business unit entity
type BusinessUnit struct {
	ID        int64
	Name      string
	Shortname string
	CompanyID int64
	LeaderID  *int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Company represents a company entity
type Company struct {
	ID   int64
	Name string
}

// BusinessUnitWithDetails represents a business unit with its leader and company
type BusinessUnitWithDetails struct {
	*BusinessUnit
	Leader  *Leader
	Company *Company
}
