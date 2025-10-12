package domain

// EnrichedData holds enriched data for starters (departments, line managers, business units)
// This is a Domain model to avoid importing Interface DTOs
type EnrichedData struct {
	Departments   map[int64]*DepartmentNested
	LineManagers  map[int64]*LineManagerNested
	BusinessUnits map[int64]*BusinessUnitNested
}

// DepartmentNested represents department with group department
type DepartmentNested struct {
	ID              int64                  `json:"id"`
	Name            string                 `json:"name"`
	Shortname       string                 `json:"shortname"`
	GroupDepartment *GroupDepartmentNested `json:"group_department,omitempty"`
}

// GroupDepartmentNested represents a parent department
type GroupDepartmentNested struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Shortname string `json:"shortname"`
}

// LineManagerNested represents a line manager
type LineManagerNested struct {
	ID       int64  `json:"id"`
	Domain   string `json:"domain"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	JobTitle string `json:"job_title"`
}

// BusinessUnitNested represents a business unit
type BusinessUnitNested struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Shortname string `json:"shortname"`
}
