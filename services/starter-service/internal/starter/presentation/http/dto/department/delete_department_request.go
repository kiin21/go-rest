package department

// DeleteDepartmentRequest holds the URI parameters for department deletion.
type DeleteDepartmentRequest struct {
	ID int64 `uri:"id" binding:"required,gt=0"`
}
