package http

import (
	"errors"
	"net/http"

	"github.com/kiin21/go-rest/pkg/httputil"
	sharedDomain "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/error"
)

func mapServiceError(err error, notFoundMsg, genericMsg string) *httputil.APIError {
	if err == nil {
		return nil
	}

	var apiErr *httputil.APIError
	if errors.As(err, &apiErr) {
		return apiErr
	}

	if notFoundMsg != "" && errors.Is(err, sharedDomain.ErrNotFound) {
		return httputil.NewAPIError(http.StatusNotFound, notFoundMsg, err.Error())
	}

	message := genericMsg
	if message == "" {
		message = "Internal server error"
	}
	return httputil.NewAPIError(http.StatusInternalServerError, message, err.Error())
}
