package dto

import (
	"errors"
	"fmt"
)

// LeaderInfo represents the leader information in the request
type LeaderInfo struct {
	ID     *int64  `json:"id" binding:"omitempty,gt=0"`
	Domain *string `json:"domain" binding:"omitempty,min=1"`
}

// AssignLeaderRequest represents the HTTP request for assigning a leader to a department
type AssignLeaderRequest struct {
	Leader LeaderInfo `json:"leader" binding:"required"`
}

// Validate ensures that either ID or Domain is provided, but not both
func (r *AssignLeaderRequest) Validate() error {
	hasID := r.Leader.ID != nil
	hasDomain := r.Leader.Domain != nil

	if !hasID && !hasDomain {
		return errors.New("either leader.id or leader.domain must be provided")
	}

	if hasID && hasDomain {
		return errors.New("cannot provide both leader.id and leader.domain")
	}

	return nil
}

// GetLeaderIdentifier returns the leader identifier and its type
func (r *AssignLeaderRequest) GetLeaderIdentifier() (interface{}, string, error) {
	if err := r.Validate(); err != nil {
		return nil, "", err
	}

	if r.Leader.ID != nil {
		return *r.Leader.ID, "id", nil
	}

	if r.Leader.Domain != nil {
		return *r.Leader.Domain, "domain", nil
	}

	return nil, "", fmt.Errorf("no valid leader identifier found")
}
