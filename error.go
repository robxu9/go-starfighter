package starfighter

import "fmt"

// APIError is for when the request processes, but returns ok = false.
// The message set is the one returned in the JSON response.
type APIError struct {
	Code    int
	Message string
}

// Error is the error string
func (a *APIError) Error() string {
	return fmt.Sprintf("starfighter api error (%d): %s", a.Code, a.Message)
}
