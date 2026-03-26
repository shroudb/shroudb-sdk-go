// ShroudbKeep error types.
//
// Auto-generated from shroudb-keep protocol spec. Do not edit.

package shroudb_keep

import (
	"fmt"
	"strings"
)

// ShroudbKeepError represents an error returned by the ShroudbKeep server.
type ShroudbKeepError struct {
	// Machine-readable error code (e.g. "NOTFOUND", "DENIED").
	Code string

	// Human-readable error message.
	Message string
}

func (e *ShroudbKeepError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func parseError(payload string) *ShroudbKeepError {
	code, message, _ := strings.Cut(payload, " ")
	return &ShroudbKeepError{Code: code, Message: message}
}

// Error code constants.
const (
	// ErrBadarg — Missing or invalid argument
	ErrBadarg = "BADARG"
	// ErrDeleted — Secret has been soft-deleted
	ErrDeleted = "DELETED"
	// ErrDenied — Authentication required or insufficient permissions
	ErrDenied = "DENIED"
	// ErrInternal — Unexpected server error
	ErrInternal = "INTERNAL"
	// ErrNotfound — Secret path or version not found
	ErrNotfound = "NOTFOUND"
	// ErrNotready — Server is starting up or shutting down
	ErrNotready = "NOTREADY"
	// ErrStorage — Backend storage error
	ErrStorage = "STORAGE"
)

// IsCode reports whether the error has the given code.
func IsCode(err error, code string) bool {
	if ke, ok := err.(*ShroudbKeepError); ok {
		return ke.Code == code
	}
	return false
}
