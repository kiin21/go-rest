package valueobject

import (
	"testing"

	domainErr "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/error"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError error
		expectValue string
	}{
		{
			name:        "valid email",
			input:       "test@vng.com.vn",
			expectError: nil,
			expectValue: "test@vng.com.vn",
		},
		{
			name:        "valid email with uppercase",
			input:       "Test@VNG.COM.VN",
			expectError: nil,
			expectValue: "Test@VNG.COM.VN",
		},
		{
			name:        "valid email with spaces",
			input:       "  test@vng.com.vn  ",
			expectError: nil,
			expectValue: "test@vng.com.vn",
		},
		{
			name:        "empty email",
			input:       "",
			expectError: domainErr.ErrEmailRequired,
		},
		{
			name:        "whitespace only",
			input:       "   ",
			expectError: domainErr.ErrEmailRequired,
		},
		{
			name:        "invalid domain",
			input:       "test@gmail.com",
			expectError: domainErr.ErrEmailInvalidDomain,
		},
		{
			name:        "invalid domain - wrong suffix",
			input:       "test@vng.com",
			expectError: domainErr.ErrEmailInvalidDomain,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := NewEmail(tt.input)

			if tt.expectError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectError)
					return
				}
				if err != tt.expectError {
					t.Errorf("expected error %v, got %v", tt.expectError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if email.Value() != tt.expectValue {
				t.Errorf("expected value %s, got %s", tt.expectValue, email.Value())
			}

			if email.String() != tt.expectValue {
				t.Errorf("expected string %s, got %s", tt.expectValue, email.String())
			}
		})
	}
}

