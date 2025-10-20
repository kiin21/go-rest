package valueobject

import (
	"strings"

	domainErr "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/error"
)

const allowedEmailDomain = "@vng.com.vn"

type Email struct {
	value string
}

func NewEmail(raw string) (Email, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return Email{}, domainErr.ErrEmailRequired
	}

	lower := strings.ToLower(trimmed)
	if !strings.HasSuffix(lower, allowedEmailDomain) {
		return Email{}, domainErr.ErrEmailInvalidDomain
	}

	return Email{value: trimmed}, nil
}

func (e Email) Value() string {
	return e.value
}

func (e Email) String() string {
	return e.value
}
