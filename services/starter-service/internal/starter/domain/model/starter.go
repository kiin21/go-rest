package model

import (
	"errors"
	"time"

	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/valueobject"
)

type StarterListFilter struct {
	DepartmentID   *int64
	BusinessUnitID *int64
	LineManagerID  *int64
	SortBy         string
	SortOrder      string
}

type Starter struct {
	ID            int64
	Domain        string
	Name          string
	Email         valueobject.Email
	Mobile        string
	WorkPhone     string
	JobTitle      string
	DepartmentID  *int64
	LineManagerID *int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
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

	emailVO, err := valueobject.NewEmail(email)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Starter{
		Domain:        domain,
		Name:          name,
		Email:         emailVO,
		Mobile:        mobile,
		WorkPhone:     workPhone,
		JobTitle:      jobTitle,
		DepartmentID:  departmentID,
		LineManagerID: lineManagerID,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func Rehydrate(
	id int64,
	domain string,
	name string,
	email string,
	mobile string,
	workPhone string,
	jobTitle string,
	departmentID *int64,
	lineManagerID *int64,
	createdAt time.Time,
	updatedAt time.Time,
) (*Starter, error) {
	emailVO, err := valueobject.NewEmail(email)
	if err != nil {
		return nil, err
	}

	return &Starter{
		ID:            id,
		Domain:        domain,
		Name:          name,
		Email:         emailVO,
		Mobile:        mobile,
		WorkPhone:     workPhone,
		JobTitle:      jobTitle,
		DepartmentID:  departmentID,
		LineManagerID: lineManagerID,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}, nil
}

// Email returns the email value as a string
func (s *Starter) GetEmail() string { return s.Email.Value() }

func (s *Starter) UpdateInfo(domain, name, email, mobile, workPhone, jobTitle string, departmentID, lineManagerID *int64) error {
	emailVO, err := valueobject.NewEmail(email)
	if err != nil {
		return err
	}

	s.Domain = domain
	s.Name = name
	s.Email = emailVO
	s.Mobile = mobile
	s.WorkPhone = workPhone
	s.JobTitle = jobTitle
	s.DepartmentID = departmentID
	s.LineManagerID = lineManagerID
	s.UpdatedAt = time.Now()

	return nil
}

type StarterESDoc struct {
	id       int64
	domain   string
	name     string
	deptName string
	buName   string
}

func (s *StarterESDoc) ID() int64                { return s.id }
func (s *StarterESDoc) Domain() string           { return s.domain }
func (s *StarterESDoc) Name() string             { return s.name }
func (s *StarterESDoc) DepartmentName() string   { return s.deptName }
func (s *StarterESDoc) BusinessUnitName() string { return s.buName }

func NewStarterESDocFromStarter(starter *Starter, enriched *EnrichedData) *StarterESDoc {
	if starter == nil {
		return nil
	}

	var deptName, buName string

	if enriched != nil {
		if depIDPtr := starter.DepartmentID; depIDPtr != nil {
			depID := *depIDPtr

			// Departments map theo department_id
			if dep, ok := enriched.Departments[depID]; ok && dep != nil {
				deptName = dep.Name

				if bu, ok := enriched.BusinessUnits[dep.ID]; ok && bu != nil {
					buName = bu.Name
				}
			}
		}
	}

	return &StarterESDoc{
		id:       starter.ID,
		domain:   starter.Domain,
		name:     starter.Name,
		deptName: deptName,
		buName:   buName,
	}
}
