package dto

// CompanyNested represents a nested company object in a response.
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

// LeaderNested represents a nested leader object in a response.
type LeaderNested struct {
	ID       int64  `json:"id"`
	Domain   string `json:"domain"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	JobTitle string `json:"job_title"`
}
