package httputil

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CustomValidator interface {
	Validate() error
}

func ValidateReq(ctx *gin.Context, req interface{}) error {
	// Bind JSON
	if err := ctx.ShouldBindJSON(req); err != nil {
		return handleValidationError(err, "Invalid request body")
	}

	// Check if req implements the CustomValidator interface
	if v, ok := req.(CustomValidator); ok {
		if err := v.Validate(); err != nil {
			return NewAPIError(http.StatusBadRequest, "Validation failed", err.Error())
		}
	}

	return nil
}

func ValidateURI(ctx *gin.Context, req interface{}) error {
	if err := ctx.ShouldBindUri(req); err != nil {
		return handleValidationError(err, "Invalid request param")
	}
	return nil
}

func ValidateQuery(ctx *gin.Context, req interface{}) error {
	if err := ctx.ShouldBindQuery(req); err != nil {
		return handleValidationError(err, "Invalid query param")
	}
	return nil
}

func handleValidationError(err error, defaultMessage string) error {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		fe := ve[0]
		errorDetail := map[string]string{
			"field":   fe.Field(),
			"message": getValidationMessage(fe),
		}
		return NewAPIError(http.StatusBadRequest, "Validation failed", errorDetail)
	}

	return NewAPIError(http.StatusBadRequest, defaultMessage, err.Error())
}

func getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s", fe.Field(), fe.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email", fe.Field())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", fe.Field(), fe.Param())
	default:
		return fe.Error()
	}
}
