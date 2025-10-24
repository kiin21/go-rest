package starter

import (
	"time"

	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/presentation/http/dto/shared"
)

type StarterResponse struct {
	ID           int64                      `json:"id"`
	Domain       string                     `json:"domain"`
	Name         string                     `json:"name"`
	Email        string                     `json:"email"`
	Mobile       string                     `json:"mobile"`
	WorkPhone    string                     `json:"work_phone"`
	JobTitle     string                     `json:"job_title"`
	Department   *shared.DepartmentNested   `json:"department,omitempty"`
	LineManager  *shared.LineManagerNested  `json:"line_manager,omitempty"`
	BusinessUnit *shared.BusinessUnitNested `json:"business_unit,omitempty"`
	CreatedAt    time.Time                  `json:"created_at"`
	UpdatedAt    time.Time                  `json:"updated_at"`
}

// EnrichedData holds related data for enrichment
type EnrichedData struct {
	Departments   map[int64]*shared.DepartmentNested
	LineManagers  map[int64]*shared.LineManagerNested
	BusinessUnits map[int64]*shared.BusinessUnitNested
}

// FromDomainEnrichment adapts domain enrichment data to HTTP DTO structures.
func FromDomainEnrichment(enriched *model.EnrichedData) *EnrichedData {
	result := &EnrichedData{
		Departments:   make(map[int64]*shared.DepartmentNested),
		LineManagers:  make(map[int64]*shared.LineManagerNested),
		BusinessUnits: make(map[int64]*shared.BusinessUnitNested),
	}

	if enriched == nil {
		return result
	}

	copyDepartment := func(src *model.DepartmentNested) *shared.DepartmentNested {
		if src == nil {
			return nil
		}
		dst := &shared.DepartmentNested{
			ID:        src.ID,
			Name:      src.Name,
			Shortname: src.Shortname,
		}
		if src.GroupDepartment != nil {
			dst.GroupDepartment = &shared.GroupDepartmentNested{
				ID:        src.GroupDepartment.ID,
				Name:      src.GroupDepartment.Name,
				Shortname: src.GroupDepartment.Shortname,
			}
		}
		return dst
	}

	for id, dept := range enriched.Departments {
		result.Departments[id] = copyDepartment(dept)
	}

	for id, manager := range enriched.LineManagers {
		if manager == nil {
			continue
		}
		result.LineManagers[id] = &shared.LineManagerNested{
			ID:       manager.ID,
			Domain:   manager.Domain,
			Name:     manager.Name,
			Email:    manager.Email,
			JobTitle: manager.JobTitle,
		}
	}

	for id, bu := range enriched.BusinessUnits {
		if bu == nil {
			continue
		}
		result.BusinessUnits[id] = &shared.BusinessUnitNested{
			ID:        bu.ID,
			Name:      bu.Name,
			Shortname: bu.Shortname,
		}
	}

	return result
}

// FromDomainEnriched converts domain document to enriched response with related data
func FromDomainEnriched(starter *model.Starter, enriched *EnrichedData) *StarterResponse {
	if starter == nil {
		return nil
	}

	response := &StarterResponse{
		ID:        starter.ID,
		Domain:    starter.Domain,
		Name:      starter.Name,
		Email:     starter.GetEmail(),
		Mobile:    starter.Mobile,
		WorkPhone: starter.WorkPhone,
		JobTitle:  starter.JobTitle,
		CreatedAt: starter.CreatedAt,
		UpdatedAt: starter.UpdatedAt,
	}

	// Map DepartmentName
	if starter.DepartmentID != nil && enriched != nil {
		deptID := *starter.DepartmentID
		if dept, ok := enriched.Departments[deptID]; ok {
			response.Department = dept

			if bu, ok := enriched.BusinessUnits[deptID]; ok {
				response.BusinessUnit = bu
			}
		}
	} // Add LM if exists
	if starter.LineManagerID != nil && enriched != nil && enriched.LineManagers != nil {
		if manager, ok := enriched.LineManagers[*starter.LineManagerID]; ok {
			response.LineManager = manager
		}
	}

	return response
}

// FromStartersEnriched converts multiple domain entities to enriched responses
func FromStartersEnriched(starters []*model.Starter, enriched *EnrichedData) []*StarterResponse {
	responses := make([]*StarterResponse, len(starters))
	for i, starter := range starters {
		responses[i] = FromDomainEnriched(starter, enriched)
	}
	return responses
}
