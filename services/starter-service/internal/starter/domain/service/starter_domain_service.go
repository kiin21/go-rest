package service

import (
	"context"

	sharedDomain "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/error"
	repo "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository"
)

type StarterDomainService struct {
	repo repo.StarterRepository
}

func NewStarterDomainService(repo repo.StarterRepository) *StarterDomainService {
	return &StarterDomainService{
		repo: repo,
	}
}

func (s *StarterDomainService) IsDomainAvailable(ctx context.Context, domain string) (bool, error) {
	starter, err := s.repo.FindByDomain(ctx, domain)
	if err != nil {
		if err.Error() == "resource not found" {
			return true, nil
		}
		return false, err
	}

	if starter != nil {
		return false, nil
	}

	return true, nil
}

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
