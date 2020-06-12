package createsend

import (
	"errors"
	"testing"
)

func TestError_Error(t *testing.T) {
	testCases := []struct {
		title    string
		expected string
		internal error
	}{
		{
			title:    "wrapped error",
			internal: errors.New("internal"),
			expected: "-1: msg. internal",
		},
		{
			title:    "no internal error",
			expected: "-1: " + ErrCodeDataProcessing.String(),
		},
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {
			err := newClientError(ErrCodeDataProcessing)
			if tC.internal != nil {
				err = newWrappedClientError("msg", tC.internal, ErrCodeDataProcessing)
			}

			var csErr *Error
			if !errors.As(err, &csErr) {
				t.Fatalf("Expected custom createsend error type")
			}

			actual := csErr.Error()
			if actual != tC.expected {
				t.Errorf("Expected: '%s', Actual: %s", tC.expected, actual)
			}
		})
	}
}
