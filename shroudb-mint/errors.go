// ShroudbMint error types.
//
// Auto-generated from shroudb-mint protocol spec. Do not edit.

package shroudb_mint

import (
	"fmt"
	"strings"
)

// ShroudbMintError represents an error returned by the ShroudbMint server.
type ShroudbMintError struct {
	// Machine-readable error code (e.g. "NOTFOUND", "DENIED").
	Code string

	// Human-readable error message.
	Message string
}

func (e *ShroudbMintError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func parseError(payload string) *ShroudbMintError {
	code, message, _ := strings.Cut(payload, " ")
	return &ShroudbMintError{Code: code, Message: message}
}

// Error code constants.
const (
	// ErrBadarg — Missing or invalid argument
	ErrBadarg = "BADARG"
	// ErrDenied — Authentication required or insufficient permissions
	ErrDenied = "DENIED"
	// ErrDisabled — CA is disabled
	ErrDisabled = "DISABLED"
	// ErrExists — CA already exists
	ErrExists = "EXISTS"
	// ErrInternal — Unexpected server error
	ErrInternal = "INTERNAL"
	// ErrNokey — Signing key not available
	ErrNokey = "NOKEY"
	// ErrNotfound — CA or certificate not found
	ErrNotfound = "NOTFOUND"
	// ErrNotready — Server is starting up or shutting down
	ErrNotready = "NOTREADY"
	// ErrStorage — Backend storage error
	ErrStorage = "STORAGE"
)

// IsCode reports whether the error has the given code.
func IsCode(err error, code string) bool {
	if ke, ok := err.(*ShroudbMintError); ok {
		return ke.Code == code
	}
	return false
}
