package createsend

import (
	"errors"
	"fmt"
)

// Error wraps server side and client side errors.
type Error struct {
	// Code error code.
	Code int
	// Message error message.
	Message string
	err     error
	wrapped bool
}

func newWrappedClientError(msg string, err error, code ClientErrorCode) error {
	return &Error{
		Code:    int(code),
		Message: msg,
		err:     err,
		wrapped: true,
	}
}

func newClientError(code ClientErrorCode) error {
	msg := code.String()
	return &Error{
		Code:    int(code),
		Message: msg,
		err:     errors.New(msg),
	}
}

// IsFromServer returns true if the error reported by the server.
//
// The method will return false if this is a client side error.
func (e *Error) IsFromServer() bool {
	return e.Code >= 0
}

// Unwrap returns the internal error.
func (e *Error) Unwrap() error {
	return e.err
}

// Error returns the string representation of the error.
func (e *Error) Error() string {
	if e.wrapped {
		return fmt.Sprintf("%d: %s. %v", e.Code, e.Message, e.err)
	}
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// Is returns true if the error is of the same type as the target.
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.Code == t.Code
}
