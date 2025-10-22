package events

type LeaderAssignmentEventPayload struct {
	FromStarter string `json:"from_starter"`
	ToStarter   string `json:"to_starter"`
	Message     string `json:"message"`
}
