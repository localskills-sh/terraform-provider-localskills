package client

import (
	"errors"
	"fmt"
)

type ApiError struct {
	StatusCode int
	Message    string
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
}

func IsNotFound(err error) bool {
	var apiErr *ApiError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 404
	}
	return false
}

func IsUnauthorized(err error) bool {
	var apiErr *ApiError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 401 || apiErr.StatusCode == 403
	}
	return false
}
