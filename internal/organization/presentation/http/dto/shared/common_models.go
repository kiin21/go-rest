package shared

// CompanyNested represents a nested company object in a response payload.
type CompanyNested struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// BusinessUnitNested represents a nested business unit object in a response.
type BusinessUnitNested struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Shortname string `json:"shortname,omitempty"`
}

// LineManagerNested represents a nested leader or line manager object in a response.
type LineManagerNested struct {
	ID       int64  `json:"id"`
	Domain   string `json:"domain"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	JobTitle string `json:"job_title"`
}

// GroupDepartmentNested represents a nested group department object in a response.
type GroupDepartmentNested struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Shortname string `json:"shortname"`
}

// DepartmentNested represents a nested department object in a response.
type DepartmentNested struct {
	ID              int64                  `json:"id"`
	Name            string                 `json:"name"`
	Shortname       string                 `json:"shortname"`
	GroupDepartment *GroupDepartmentNested `json:"group_department,omitempty"`
}

// OrgDepartmentNested represents a nested department within the organization context.
type OrgDepartmentNested struct {
	ID        int64  `json:"id"`
	FullName  string `json:"full_name"`
	Shortname string `json:"shortname"`
}
