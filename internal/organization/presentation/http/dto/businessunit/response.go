package businessunit

import (
	"time"

	model "github.com/kiin21/go-rest/internal/organization/domain/model"
	"github.com/kiin21/go-rest/internal/organization/presentation/http/dto/shared"
)

// BusinessUnitResponse represents business unit payloads returned in API.
type BusinessUnitResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Shortname string    `json:"shortname"`
	CompanyID int64     `json:"company_id"`
	LeaderID  *int64    `json:"leader_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BusinessUnitDetailResponse represents business unit with nested details.
type BusinessUnitDetailResponse struct {
	ID        int64                     `json:"id"`
	Name      string                    `json:"name"`
	Shortname string                    `json:"shortname"`
	Company   *shared.CompanyNested     `json:"company,omitempty"`
	Leader    *shared.LineManagerNested `json:"leader,omitempty"`
	CreatedAt time.Time                 `json:"created_at"`
	UpdatedAt time.Time                 `json:"updated_at"`
}

// FromBusinessUnitWithDetails converts a detailed domain entity to a detailed response DTO.
func FromBusinessUnitWithDetails(unit *model.BusinessUnitWithDetails) *BusinessUnitDetailResponse {
	if unit == nil {
		return nil
	}

	resp := &BusinessUnitDetailResponse{
		ID:        unit.ID,
		Name:      unit.Name,
		Shortname: unit.Shortname,
		CreatedAt: unit.CreatedAt,
		UpdatedAt: unit.UpdatedAt,
	}

	if unit.Company != nil {
		resp.Company = &shared.CompanyNested{
			ID:   unit.Company.ID,
			Name: unit.Company.Name,
		}
	}

	if unit.Leader != nil {
		resp.Leader = &shared.LineManagerNested{
			ID:       unit.Leader.ID,
			Domain:   unit.Leader.Domain,
			Name:     unit.Leader.Name,
			Email:    unit.Leader.Email,
			JobTitle: unit.Leader.JobTitle,
		}
	}

	return resp
}

// FromBusinessUnitsWithDetails converts a slice of detailed business units to DTOs.
func FromBusinessUnitsWithDetails(units []*model.BusinessUnitWithDetails) []*BusinessUnitDetailResponse {
	responses := make([]*BusinessUnitDetailResponse, len(units))
	for i, unit := range units {
		responses[i] = FromBusinessUnitWithDetails(unit)
	}
	return responses
}
