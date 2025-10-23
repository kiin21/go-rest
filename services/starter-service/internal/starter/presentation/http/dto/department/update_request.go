package department

import "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/department/command"

type UpdateDepartmentRequest struct {
	FullName          *string `json:"full_name" binding:"omitempty,min=3,max=255"`
	Shortname         *string `json:"shortname" binding:"omitempty,min=2,max=100"`
	BusinessUnitID    *int64  `json:"business_unit_id" binding:"omitempty,gt=0"`
	GroupDepartmentID *int64  `json:"group_department_id" binding:"omitempty,gt=0"`
	LeaderID          *int64  `json:"leader_id" binding:"omitempty,gt=0"`
}

func (r *UpdateDepartmentRequest) ToCommand(deptId int64) *command.UpdateDepartmentCommand {
	return &command.UpdateDepartmentCommand{
		ID:                deptId,
		FullName:          r.FullName,
		Shortname:         r.Shortname,
		BusinessUnitID:    r.BusinessUnitID,
		GroupDepartmentID: r.GroupDepartmentID,
		LeaderID:          r.LeaderID,
	}
}
