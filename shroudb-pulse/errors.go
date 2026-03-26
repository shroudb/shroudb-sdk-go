// ShroudbPulse error types.
//
// Auto-generated from shroudb-pulse protocol spec. Do not edit.

package shroudb_pulse

import (
	"fmt"
	"strings"
)

// ShroudbPulseError represents an error returned by the ShroudbPulse server.
type ShroudbPulseError struct {
	// Machine-readable error code (e.g. "NOTFOUND", "DENIED").
	Code string

	// Human-readable error message.
	Message string
}

func (e *ShroudbPulseError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func parseError(payload string) *ShroudbPulseError {
	code, message, _ := strings.Cut(payload, " ")
	return &ShroudbPulseError{Code: code, Message: message}
}

// Error code constants.
const (
	// ErrBadarg — Missing or invalid argument
	ErrBadarg = "BADARG"
	// ErrDenied — Authentication required or insufficient permissions
	ErrDenied = "DENIED"
	// ErrInternal — Unexpected server error
	ErrInternal = "INTERNAL"
	// ErrNotfound — Resource not found
	ErrNotfound = "NOTFOUND"
	// ErrNotready — Server is starting up or shutting down
	ErrNotready = "NOTREADY"
	// ErrStorage — Backend storage or WAL error
	ErrStorage = "STORAGE"
)

// IsCode reports whether the error has the given code.
func IsCode(err error, code string) bool {
	if ke, ok := err.(*ShroudbPulseError); ok {
		return ke.Code == code
	}
	return false
}
