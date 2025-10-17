package dto

import (
	"time"

	shareddto "github.com/kiin21/go-rest/internal/shared/dto"
	starterAggregate "github.com/kiin21/go-rest/internal/starter/domain/aggregate"
)

type DepartmentNested = shareddto.DepartmentNested
type GroupDepartmentNested = shareddto.GroupDepartmentNested
type LineManagerNested = shareddto.LineManagerNested
type BusinessUnitNested = shareddto.BusinessUnitNested

type StarterResponse struct {
	ID           int64               `json:"id"`
	Domain       string              `json:"domain"`
	Name         string              `json:"name"`
	Email        string              `json:"email"`
	Mobile       string              `json:"mobile"`
	WorkPhone    string              `json:"work_phone"`
	JobTitle     string              `json:"job_title"`
	Department   *DepartmentNested   `json:"department,omitempty"`
	LineManager  *LineManagerNested  `json:"line_manager,omitempty"`
	BusinessUnit *BusinessUnitNested `json:"business_unit,omitempty"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

// StarterEnrichedData holds related data for enrichment
type StarterEnrichedData struct {
	Departments   map[int64]*DepartmentNested
	LineManagers  map[int64]*LineManagerNested
	BusinessUnits map[int64]*BusinessUnitNested
}

// FromDomainEnriched converts domain entity to enriched response with related data
func FromDomainEnriched(starter *starterAggregate.Starter, enriched *StarterEnrichedData) *StarterResponse {
	if starter == nil {
		return nil
	}

	response := &StarterResponse{
		ID:        starter.ID(),
		Domain:    starter.Domain(),
		Name:      starter.Name(),
		Email:     starter.Email(),
		Mobile:    starter.Mobile(),
		WorkPhone: starter.WorkPhone(),
		JobTitle:  starter.JobTitle(),
		CreatedAt: starter.CreatedAt(),
		UpdatedAt: starter.UpdatedAt(),
	}

	// Map Department
	if starter.DepartmentID() != nil && enriched != nil {
		deptID := *starter.DepartmentID()
		if dept, ok := enriched.Departments[deptID]; ok {
			response.Department = dept

			// Map BusinessUnit using department ID as key
			if bu, ok := enriched.BusinessUnits[deptID]; ok {
				response.BusinessUnit = bu
			}
		}
	} // Add LM if exists
	if starter.LineManagerID() != nil && enriched != nil && enriched.LineManagers != nil {
		if manager, ok := enriched.LineManagers[*starter.LineManagerID()]; ok {
			response.LineManager = manager
		}
	}

	return response
}

// FromStartersEnriched converts multiple domain entities to enriched responses
func FromStartersEnriched(starters []*starterAggregate.Starter, enriched *StarterEnrichedData) []*StarterResponse {
	responses := make([]*StarterResponse, len(starters))
	for i, starter := range starters {
		responses[i] = FromDomainEnriched(starter, enriched)
	}
	return responses
}
