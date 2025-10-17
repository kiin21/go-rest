package domain

type BusinessUnitNested struct {
	ID        int64
	Name      string
	Shortname string
}

type LineManagerNested struct {
	ID       int64
	Domain   string
	Name     string
	Email    string
	JobTitle string
}

type GroupDepartmentNested struct {
	ID        int64
	Name      string
	Shortname string
}

type DepartmentNested struct {
	ID              int64
	Name            string
	Shortname       string
	GroupDepartment *GroupDepartmentNested
}

type OrgDepartmentNested struct {
	ID        int64
	FullName  string
	Shortname string
}
