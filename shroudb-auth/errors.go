// ShroudbAuth error types.
//
// Auto-generated from shroudb-auth protocol spec. Do not edit.

package shroudb_auth

import (
	"fmt"
)

// ShroudbAuthError represents an error returned by the ShroudbAuth server.
type ShroudbAuthError struct {
	// Machine-readable error code (e.g. "UNAUTHORIZED", "CONFLICT").
	Code string `json:"code"`

	// Human-readable error message.
	Message string `json:"message"`
}

func (e *ShroudbAuthError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Error code constants.
const (
	// ErrBadRequest — Invalid request body or parameters
	ErrBadRequest = "BAD_REQUEST"
	// ErrConflict — Resource already exists (e.g. duplicate signup)
	ErrConflict = "CONFLICT"
	// ErrForbidden — Insufficient permissions
	ErrForbidden = "FORBIDDEN"
	// ErrInternal — Internal server error
	ErrInternal = "INTERNAL"
	// ErrTooManyRequests — Account locked due to too many failed attempts
	ErrTooManyRequests = "TOO_MANY_REQUESTS"
	// ErrUnauthorized — Authentication required or invalid credentials
	ErrUnauthorized = "UNAUTHORIZED"
)

// IsCode reports whether the error has the given error code.
func IsCode(err error, code string) bool {
	if ae, ok := err.(*ShroudbAuthError); ok {
		return ae.Code == code
	}
	return false
}


// IsBadRequest reports whether the error is BAD_REQUEST: Invalid request body or parameters
func IsBadRequest(err error) bool { return IsCode(err, ErrBadRequest) }

// IsConflict reports whether the error is CONFLICT: Resource already exists (e.g. duplicate signup)
func IsConflict(err error) bool { return IsCode(err, ErrConflict) }

// IsForbidden reports whether the error is FORBIDDEN: Insufficient permissions
func IsForbidden(err error) bool { return IsCode(err, ErrForbidden) }

// IsInternal reports whether the error is INTERNAL: Internal server error
func IsInternal(err error) bool { return IsCode(err, ErrInternal) }

// IsTooManyRequests reports whether the error is TOO_MANY_REQUESTS: Account locked due to too many failed attempts
func IsTooManyRequests(err error) bool { return IsCode(err, ErrTooManyRequests) }

// IsUnauthorized reports whether the error is UNAUTHORIZED: Authentication required or invalid credentials
func IsUnauthorized(err error) bool { return IsCode(err, ErrUnauthorized) }
