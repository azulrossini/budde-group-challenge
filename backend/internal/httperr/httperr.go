package httperr

import (
	"encoding/json"
	"net/http"

	"github.com/azulrossini/budde-group-challenge/backend/internal/api"
	"github.com/azulrossini/budde-group-challenge/backend/internal/apperr"
)

const (
	headerContentType = "Content-Type"
	contentTypeJSON   = "application/json"
)

// Body maps an apperr.Error to the HTTP status and ErrorResponse body it
// should produce, per docs/SPEC.md §8.
func Body(err *apperr.Error, requestID string) (int, api.ErrorResponse) {
	body := api.ErrorResponse{
		Code:      api.ErrorCode(err.Code),
		Message:   err.Message,
		RequestId: requestID,
	}

	if len(err.Details) > 0 {
		details := make([]api.FieldError, len(err.Details))
		for i, d := range err.Details {
			details[i] = api.FieldError{Field: d.Field, Message: d.Message}
		}
		body.Details = &details
	}

	return statusFor(err.Code), body
}

// Write writes err directly to w. Used outside the generated strict
// handler's typed responses: middleware, and the wrapper's own error hooks.
func Write(w http.ResponseWriter, err *apperr.Error, requestID string) {
	status, body := Body(err, requestID)
	w.Header().Set(headerContentType, contentTypeJSON)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func statusFor(code apperr.Code) int {
	switch code {
	case apperr.CodeValidationError:
		return http.StatusUnprocessableEntity
	case apperr.CodeBadRequest:
		return http.StatusBadRequest
	case apperr.CodeNotFound:
		return http.StatusNotFound
	case apperr.CodeUnauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
