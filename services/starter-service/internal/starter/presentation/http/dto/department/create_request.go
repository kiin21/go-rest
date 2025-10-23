package department

import "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/department/command"

type CreateDepartmentRequest struct {
	FullName          string `json:"full_name" binding:"required,min=3,max=255"`
	Shortname         string `json:"shortname" binding:"required,min=2,max=100"`
	BusinessUnitID    *int64 `json:"business_unit_id" binding:"omitempty,gt=0"`
	GroupDepartmentID *int64 `json:"group_department_id" binding:"omitempty,gt=0"`
	LeaderID          *int64 `json:"leader_id" binding:"omitempty,gt=0"`
}

func (r *CreateDepartmentRequest) ToCommand() *command.CreateDepartmentCommand {
	return &command.CreateDepartmentCommand{
		FullName:          r.FullName,
		Shortname:         r.Shortname,
		BusinessUnitID:    r.BusinessUnitID,
		GroupDepartmentID: r.GroupDepartmentID,
		LeaderID:          r.LeaderID,
	}
}
