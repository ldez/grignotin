package goproxy

import "fmt"

// APIError represents an error from the GoProxy.
type APIError struct {
	StatusCode int
	Message    string
}

func (a *APIError) Error() string {
	return fmt.Sprintf("error: %d: %s", a.StatusCode, a.Message)
}
