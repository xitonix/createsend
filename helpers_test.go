package createsend_test

import (
	"testing"

	"github.com/xitonix/createsend"
)

func checkErrorType(t *testing.T, err error) {
	t.Helper()
	if _, ok := err.(*createsend.Error); !ok {
		t.Error("We should always return a custom createsend Error type")
	}
}
