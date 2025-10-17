package aggregate

import (
	"errors"
	"time"
)

// Starter represents the aggregate root for employee onboarding.
type Starter struct {
	id            int64
	domain        string
	name          string
	email         string
	mobile        string
	workPhone     string
	jobTitle      string
	departmentID  *int64
	lineManagerID *int64
	createdAt     time.Time
	updatedAt     time.Time
}

// NewStarter validates input and creates a new Starter aggregate.
func NewStarter(domain, name, email, mobile, workPhone, jobTitle string, departmentID, lineManagerID *int64) (*Starter, error) {
	if domain == "" {
		return nil, errors.New("domain is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}
	if mobile == "" {
		return nil, errors.New("mobile is required")
	}
	if jobTitle == "" {
		return nil, errors.New("job title is required")
	}

	now := time.Now()
	return &Starter{
		domain:        domain,
		name:          name,
		email:         email,
		mobile:        mobile,
		workPhone:     workPhone,
		jobTitle:      jobTitle,
		departmentID:  departmentID,
		lineManagerID: lineManagerID,
		createdAt:     now,
		updatedAt:     now,
	}, nil
}

// Rehydrate reconstructs a Starter aggregate from persistence.
func Rehydrate(
	id int64,
	domain,
	name,
	email,
	mobile,
	workPhone,
	jobTitle string,
	departmentID,
	lineManagerID *int64,
	createdAt,
	updatedAt time.Time,
) *Starter {
	return &Starter{
		id:            id,
		domain:        domain,
		name:          name,
		email:         email,
		mobile:        mobile,
		workPhone:     workPhone,
		jobTitle:      jobTitle,
		departmentID:  departmentID,
		lineManagerID: lineManagerID,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
}

// ID returns the aggregate identifier.
func (s *Starter) ID() int64 { return s.id }

// Domain returns the login domain.
func (s *Starter) Domain() string { return s.domain }

// Name returns the full name.
func (s *Starter) Name() string { return s.name }

// Email returns the email address.
func (s *Starter) Email() string { return s.email }

// Mobile returns the mobile phone number.
func (s *Starter) Mobile() string { return s.mobile }

// WorkPhone returns the office phone number.
func (s *Starter) WorkPhone() string { return s.workPhone }

// JobTitle returns the current job title.
func (s *Starter) JobTitle() string { return s.jobTitle }

// DepartmentID returns the identifier of the assigned department.
func (s *Starter) DepartmentID() *int64 { return s.departmentID }

// LineManagerID returns the identifier of the line manager.
func (s *Starter) LineManagerID() *int64 { return s.lineManagerID }

// CreatedAt returns the creation timestamp.
func (s *Starter) CreatedAt() time.Time { return s.createdAt }

// UpdatedAt returns the last update timestamp.
func (s *Starter) UpdatedAt() time.Time { return s.updatedAt }

// UpdateInfo updates the mutable properties of the starter.
func (s *Starter) UpdateInfo(name, email, mobile, workPhone, jobTitle string, departmentID, lineManagerID *int64) {
	s.name = name
	s.email = email
	s.mobile = mobile
	s.workPhone = workPhone
	s.jobTitle = jobTitle
	s.departmentID = departmentID
	s.lineManagerID = lineManagerID
	s.updatedAt = time.Now()
}

// AssignToDepartment assigns the starter to a department.
func (s *Starter) AssignToDepartment(departmentID int64) {
	s.departmentID = &departmentID
	s.updatedAt = time.Now()
}

// AssignLineManager assigns a line manager to the starter.
func (s *Starter) AssignLineManager(managerID *int64) {
	s.lineManagerID = managerID
	s.updatedAt = time.Now()
}

// CanBeDeleted checks business rules that could prevent deletion.
func (s *Starter) CanBeDeleted(hasSubordinates, isDepartmentLeader, isBusinessUnitLeader bool) bool {
	return !hasSubordinates && !isDepartmentLeader && !isBusinessUnitLeader
}

// MarkAsRemoved is a placeholder to raise domain events before deletion.
func (s *Starter) MarkAsRemoved() {}
