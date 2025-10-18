package department

type GetDepartmentRequest struct {
	ID int64 `uri:"id" binding:"required,gt=0"`
}
