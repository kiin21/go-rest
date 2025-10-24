package service

import (
	"context"
	"errors"
	"testing"

	sharedDomain "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/error"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository/mocks"
)

func TestIsDomainAvailable(t *testing.T) {
	tests := []struct {
		name            string
		domain          string
		mockFindFunc    func(ctx context.Context, domain string) (*model.Starter, error)
		expectedResult  bool
		expectError     bool
	}{
		{
			name:   "domain is available - not found",
			domain: "available",
			mockFindFunc: func(ctx context.Context, domain string) (*model.Starter, error) {
				return nil, errors.New("resource not found")
			},
			expectedResult: true,
			expectError:    false,
		},
		{
			name:   "domain is not available",
			domain: "taken",
			mockFindFunc: func(ctx context.Context, domain string) (*model.Starter, error) {
				starter, _ := model.NewStarter("taken", "Test User", "test@vng.com.vn", "0123456789", "", "Developer", nil, nil)
				return starter, nil
			},
			expectedResult: false,
			expectError:    false,
		},
		{
			name:   "repository error",
			domain: "error",
			mockFindFunc: func(ctx context.Context, domain string) (*model.Starter, error) {
				return nil, errors.New("database error")
			},
			expectedResult: false,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockStarterRepository{
				FindByDomainFunc: tt.mockFindFunc,
			}

			service := NewStarterDomainService(mockRepo)
			available, err := service.IsDomainAvailable(context.Background(), tt.domain)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if available != tt.expectedResult {
				t.Errorf("expected available=%v, got %v", tt.expectedResult, available)
			}
		})
	}
}

func TestValidateDomainUniqueness(t *testing.T) {
	tests := []struct {
		name         string
		domain       string
		mockFindFunc func(ctx context.Context, domain string) (*model.Starter, error)
		expectError  error
	}{
		{
			name:   "domain is unique",
			domain: "unique",
			mockFindFunc: func(ctx context.Context, domain string) (*model.Starter, error) {
				return nil, errors.New("resource not found")
			},
			expectError: nil,
		},
		{
			name:   "domain already exists",
			domain: "taken",
			mockFindFunc: func(ctx context.Context, domain string) (*model.Starter, error) {
				starter, _ := model.NewStarter("taken", "Test User", "test@vng.com.vn", "0123456789", "", "Developer", nil, nil)
				return starter, nil
			},
			expectError: sharedDomain.ErrDomainAlreadyExists,
		},
		{
			name:   "repository error",
			domain: "error",
			mockFindFunc: func(ctx context.Context, domain string) (*model.Starter, error) {
				return nil, errors.New("database error")
			},
			expectError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockStarterRepository{
				FindByDomainFunc: tt.mockFindFunc,
			}

			service := NewStarterDomainService(mockRepo)
			err := service.ValidateDomainUniqueness(context.Background(), tt.domain)

			if tt.expectError != nil {
				if err == nil {
					t.Error("expected error but got nil")
					return
				}
				if tt.expectError == sharedDomain.ErrDomainAlreadyExists && err != sharedDomain.ErrDomainAlreadyExists {
					t.Errorf("expected error %v, got %v", tt.expectError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

