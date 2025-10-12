package domain

import "errors"

// Common domain errors
var (
	// Resource errors
	ErrNotFound            = errors.New("resource not found")
	ErrAlreadyExists       = errors.New("resource already exists")
	ErrDuplicateEntry      = errors.New("duplicate entry")
	ErrDomainAlreadyExists = errors.New("domain already exists") // Specific for starter domain uniqueness

	// Permission errors
	ErrForbidden    = errors.New("forbidden")
	ErrUnauthorized = errors.New("unauthorized")

	// Validation errors
	ErrValidation   = errors.New("validation failed")
	ErrInvalidInput = errors.New("invalid input")

	// Business logic errors
	ErrConflict     = errors.New("conflict")
	ErrPrecondition = errors.New("precondition failed")

	// Operation errors
	ErrOperationFailed = errors.New("operation failed")
	ErrTimeout         = errors.New("operation timeout")
)
