package domain

import "errors"

// Common domain errors
var (
	ErrNotFound            = errors.New("resource not found")
	ErrDuplicateEntry      = errors.New("duplicate entry")
	ErrDomainAlreadyExists = errors.New("domain already exists") // Specific for starter domain uniqueness

	ErrValidation   = errors.New("validation failed")
	ErrInvalidInput = errors.New("invalid input")
)
