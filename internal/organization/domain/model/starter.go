package model

import (
	"errors"
	"time"
)

type StarterListFilter struct {
	DepartmentID   *int64
	BusinessUnitID *int64
	LineManagerID  *int64
}

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

func (s *Starter) ID() int64 { return s.id }

func (s *Starter) Domain() string { return s.domain }

func (s *Starter) Name() string { return s.name }

func (s *Starter) Email() string { return s.email }

func (s *Starter) Mobile() string { return s.mobile }

func (s *Starter) WorkPhone() string { return s.workPhone }

func (s *Starter) JobTitle() string { return s.jobTitle }

func (s *Starter) DepartmentID() *int64 { return s.departmentID }

func (s *Starter) LineManagerID() *int64 { return s.lineManagerID }

func (s *Starter) CreatedAt() time.Time { return s.createdAt }

func (s *Starter) UpdatedAt() time.Time { return s.updatedAt }

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
