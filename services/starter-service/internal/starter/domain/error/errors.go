package error

import "errors"

// Common domain errors
var (
	ErrNotFound            = errors.New("resource not found")
	ErrDuplicateEntry      = errors.New("duplicate entry")
	ErrDomainAlreadyExists = errors.New("domain already exists")
	ErrEmailRequired       = errors.New("email is required")
	ErrEmailInvalidDomain  = errors.New("email must end with @vng.com.vn")

	ErrValidation   = errors.New("validation failed")
	ErrInvalidInput = errors.New("invalid input")
)
