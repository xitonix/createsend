package mock

import "errors"

var (
	// ErrDeliberate occurs when a unit test simulates an error.
	ErrDeliberate = errors.New("deliberate error occurred")
)