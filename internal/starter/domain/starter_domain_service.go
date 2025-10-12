package domain

import (
	"context"

	sharedDomain "github.com/kiin21/go-rest/internal/shared/domain"
)

type StarterDomainService struct {
	repo StarterRepository
}

// NewStarterDomainService creates a new domain service
func NewStarterDomainService(repo StarterRepository) *StarterDomainService {
	return &StarterDomainService{
		repo: repo,
	}
}

func (s *StarterDomainService) IsDomainAvailable(ctx context.Context, domain string) (bool, error) {

	starter, err := s.repo.FindByDomain(ctx, domain)
	if err != nil {
		// If not found, domain is available
		if err.Error() == "resource not found" { // Or check specific error type
			return true, nil
		}
		return false, err
	}

	// If found, domain is not available
	if starter != nil {
		return false, nil
	}

	return true, nil
}

// ValidateDomainUniqueness validates that domain is unique before creation
// Returns error if domain already exists
func (s *StarterDomainService) ValidateDomainUniqueness(ctx context.Context, domain string) error {
	available, err := s.IsDomainAvailable(ctx, domain)
	if err != nil {
		return err
	}

	if !available {
		return sharedDomain.ErrDomainAlreadyExists
	}

	return nil
}
