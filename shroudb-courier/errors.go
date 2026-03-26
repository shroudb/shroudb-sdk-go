// ShroudbCourier error types.
//
// Auto-generated from shroudb-courier protocol spec. Do not edit.

package shroudb_courier

import (
	"fmt"
	"strings"
)

// ShroudbCourierError represents an error returned by the ShroudbCourier server.
type ShroudbCourierError struct {
	// Machine-readable error code (e.g. "NOTFOUND", "DENIED").
	Code string

	// Human-readable error message.
	Message string
}

func (e *ShroudbCourierError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func parseError(payload string) *ShroudbCourierError {
	code, message, _ := strings.Cut(payload, " ")
	return &ShroudbCourierError{Code: code, Message: message}
}

// Error code constants.
const (
	// ErrBadarg — Missing or invalid argument
	ErrBadarg = "BADARG"
	// ErrDeliveryFailed — Notification delivery failed
	ErrDeliveryFailed = "DELIVERY_FAILED"
	// ErrDenied — Authentication required or insufficient permissions
	ErrDenied = "DENIED"
	// ErrInternal — Unexpected server error
	ErrInternal = "INTERNAL"
	// ErrNotfound — Template not found
	ErrNotfound = "NOTFOUND"
	// ErrNotready — Server is starting up or shutting down
	ErrNotready = "NOTREADY"
	// ErrTemplateError — Template rendering error
	ErrTemplateError = "TEMPLATE_ERROR"
)

// IsCode reports whether the error has the given code.
func IsCode(err error, code string) bool {
	if ke, ok := err.(*ShroudbCourierError); ok {
		return ke.Code == code
	}
	return false
}
