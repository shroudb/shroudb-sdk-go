// ShroudbSentry error types.
//
// Auto-generated from shroudb-sentry protocol spec. Do not edit.

package shroudb_sentry

import (
	"fmt"
	"strings"
)

// ShroudbSentryError represents an error returned by the ShroudbSentry server.
type ShroudbSentryError struct {
	// Machine-readable error code (e.g. "NOTFOUND", "DENIED").
	Code string

	// Human-readable error message.
	Message string
}

func (e *ShroudbSentryError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func parseError(payload string) *ShroudbSentryError {
	code, message, _ := strings.Cut(payload, " ")
	return &ShroudbSentryError{Code: code, Message: message}
}

// Error code constants.
const (
	// ErrBadarg — Missing or invalid argument
	ErrBadarg = "BADARG"
	// ErrDenied — Authentication required or insufficient permissions
	ErrDenied = "DENIED"
	// ErrInternal — Unexpected server error
	ErrInternal = "INTERNAL"
	// ErrNokey — Signing key not available
	ErrNokey = "NOKEY"
	// ErrNotfound — Policy not found
	ErrNotfound = "NOTFOUND"
	// ErrNotready — Server is starting up or shutting down
	ErrNotready = "NOTREADY"
)

// IsCode reports whether the error has the given code.
func IsCode(err error, code string) bool {
	if ke, ok := err.(*ShroudbSentryError); ok {
		return ke.Code == code
	}
	return false
}
