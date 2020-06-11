package createsend

import (
	"errors"
	"fmt"
)

type Error struct {
	Code    int
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

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Error() string {
	if e.wrapped {
		return fmt.Sprintf("%d: %s. %v", e.Code, e.Message, e.err)
	}
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.Code == t.Code
}
