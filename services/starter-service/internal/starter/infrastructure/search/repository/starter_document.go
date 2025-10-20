package repository

import "time"

// StarterDocument mirrors the Elasticsearch document structure optimized for search.
type StarterDocument struct {
	ID            int64  `json:"id"`
	Domain        string `json:"domain"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Mobile        string `json:"mobile"`
	WorkPhone     string `json:"work_phone"`
	JobTitle      string `json:"job_title"`
	DepartmentID  *int64 `json:"department_id,omitempty"`
	LineManagerID *int64 `json:"line_manager_id,omitempty"`

	// Additional fields for search optimization.
	FullText     string   `json:"full_text"`     // Combined text for full-text search.
	SearchTokens []string `json:"search_tokens"` // Tokenized fields for better matching.

	// Metadata.
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IndexedAt time.Time `json:"indexed_at"`
}
