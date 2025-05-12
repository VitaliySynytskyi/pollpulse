package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// ServiceError represents a standard error response from services
type ServiceError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

var (
	// ErrNotFound indicates a resource was not found
	ErrNotFound = errors.New("resource not found")
	
	// ErrUnauthorized indicates authentication failure
	ErrUnauthorized = errors.New("unauthorized")
	
	// ErrForbidden indicates a lack of permission
	ErrForbidden = errors.New("forbidden")
	
	// ErrBadRequest indicates an invalid request
	ErrBadRequest = errors.New("bad request")

	// ErrInternalServer indicates an internal server error
	ErrInternalServer = errors.New("internal server error")
)

// WriteError writes a standardized error response
func WriteError(w http.ResponseWriter, err error, statusCode int, details string) {
	serviceErr := ServiceError{
		Code:    statusCode,
		Message: err.Error(),
		Details: details,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(serviceErr)
}

// HandleError maps known errors to HTTP status codes and writes a response
func HandleError(w http.ResponseWriter, err error, details string) {
	switch {
	case errors.Is(err, ErrNotFound):
		WriteError(w, err, http.StatusNotFound, details)
	case errors.Is(err, ErrUnauthorized):
		WriteError(w, err, http.StatusUnauthorized, details)
	case errors.Is(err, ErrForbidden):
		WriteError(w, err, http.StatusForbidden, details)
	case errors.Is(err, ErrBadRequest):
		WriteError(w, err, http.StatusBadRequest, details)
	default:
		// Log the original error but don't expose it to the client
		fmt.Printf("Internal error: %v\n", err)
		WriteError(w, ErrInternalServer, http.StatusInternalServerError, "")
	}
}

// NewError creates a new error with formatted message
func NewError(err error, format string, args ...interface{}) error {
	return fmt.Errorf("%w: %s", err, fmt.Sprintf(format, args...))
} 