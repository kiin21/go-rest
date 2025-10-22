package events

type IndexStarterPayload struct {
	StarterID int64  `json:"starter_id"`
	Domain    string `json:"domain"`
	Name      string `json:"name"`
	DeptName  string `json:"dept_name"`
	BUName    string `json:"bu_name"`
}
