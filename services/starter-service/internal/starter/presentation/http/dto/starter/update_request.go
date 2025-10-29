package starter

import "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/starter/command"

type UpdateStarterRequest struct {
	Domain        *string `uri:"domain" binding:"required,min=2,max=25"`
	Name          *string `json:"name" binding:"omitempty,min=2,max=255"`
	Email         *string `json:"email" binding:"omitempty,email,max=100"`
	Mobile        *string `json:"mobile" binding:"omitempty,min=10,max=20"`
	WorkPhone     *string `json:"work_phone" binding:"omitempty,max=20"`
	JobTitle      *string `json:"job_title" binding:"omitempty,min=2,max=100"`
	DepartmentID  *int64  `json:"department_id" binding:"omitempty,gt=0"`
	LineManagerID *int64  `json:"line_manager_id" binding:"omitempty,gt=0"`
}

func (r *UpdateStarterRequest) ToCommand(originalDomain string) *command.UpdateStarterCommand {
	return &command.UpdateStarterCommand{
		OriginalDomain: originalDomain,
		Domain:         r.Domain,
		Name:           r.Name,
		Email:          r.Email,
		Mobile:         r.Mobile,
		WorkPhone:      r.WorkPhone,
		JobTitle:       r.JobTitle,
		DepartmentID:   r.DepartmentID,
		LineManagerID:  r.LineManagerID,
	}
}
