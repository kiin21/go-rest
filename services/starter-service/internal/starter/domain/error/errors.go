package error

import "errors"

// Common domain errors
var (
	ErrNotFound            = errors.New("resource not found")
	ErrDomainAlreadyExists = errors.New("domain already exists")
	ErrEmailRequired       = errors.New("email is required")
	ErrEmailInvalidDomain  = errors.New("email must end with @vng.com.vn")

	ErrInvalidInput = errors.New("invalid input")
)
