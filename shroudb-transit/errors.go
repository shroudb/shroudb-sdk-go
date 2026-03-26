// ShroudbTransit error types.
//
// Auto-generated from shroudb-transit protocol spec. Do not edit.

package shroudb_transit

import (
	"fmt"
	"strings"
)

// ShroudbTransitError represents an error returned by the ShroudbTransit server.
type ShroudbTransitError struct {
	// Machine-readable error code (e.g. "NOTFOUND", "DENIED").
	Code string

	// Human-readable error message.
	Message string
}

func (e *ShroudbTransitError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func parseError(payload string) *ShroudbTransitError {
	code, message, _ := strings.Cut(payload, " ")
	return &ShroudbTransitError{Code: code, Message: message}
}

// Error code constants.
const (
	// ErrBadarg — Missing or invalid argument
	ErrBadarg = "BADARG"
	// ErrDenied — Authentication required or insufficient permissions
	ErrDenied = "DENIED"
	// ErrDisabled — Keyring is disabled
	ErrDisabled = "DISABLED"
	// ErrInternal — Unexpected server error
	ErrInternal = "INTERNAL"
	// ErrNotfound — Keyring or key version not found
	ErrNotfound = "NOTFOUND"
	// ErrNotready — Server is starting up or shutting down
	ErrNotready = "NOTREADY"
	// ErrWrongtype — Operation not supported for this keyring type
	ErrWrongtype = "WRONGTYPE"
)

// IsCode reports whether the error has the given code.
func IsCode(err error, code string) bool {
	if ke, ok := err.(*ShroudbTransitError); ok {
		return ke.Code == code
	}
	return false
}
