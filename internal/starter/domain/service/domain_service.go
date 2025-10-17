package service

import (
	"context"

	sharedDomain "github.com/kiin21/go-rest/internal/shared/domain"
	"github.com/kiin21/go-rest/internal/starter/domain/port"
)

// StarterDomainService contains domain-specific operations for the Starter aggregate.
type StarterDomainService struct {
	repo port.StarterRepository
}

// NewStarterDomainService creates a new domain service.
func NewStarterDomainService(repo port.StarterRepository) *StarterDomainService {
	return &StarterDomainService{
		repo: repo,
	}
}

// IsDomainAvailable checks whether the given domain identifier is available.
func (s *StarterDomainService) IsDomainAvailable(ctx context.Context, domain string) (bool, error) {
	starter, err := s.repo.FindByDomain(ctx, domain)
	if err != nil {
		// If not found, domain is available.
		if err.Error() == "resource not found" {
			return true, nil
		}
		return false, err
	}

	// If found, domain is not available.
	if starter != nil {
		return false, nil
	}

	return true, nil
}

// ValidateDomainUniqueness ensures the domain is unique before creation.
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
