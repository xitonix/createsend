package transactional_test

import (
	"testing"

	"github.com/xitonix/createsend/transactional"
)

func TestDefaultOptions(t *testing.T) {
	ops := transactional.Options{}
	if ops.ClientID() != "" {
		t.Errorf("Expected client ID: '', Actual: %s", ops.ClientID())
	}

	if ops.SmartEmailStatus() != transactional.UnknownSmartEmail {
		t.Errorf("Expected smart email status: %s, Actual: %s", transactional.UnknownSmartEmail.String(), ops.SmartEmailStatus().String())
	}
}

func TestWithClientID(t *testing.T) {
	const expected = "client_id"
	ops := &transactional.Options{}
	option := transactional.WithClientID(expected)
	option(ops)
	actual := ops.ClientID()
	if actual != expected {
		t.Errorf("Expected client ID: %s, Actual: %s", expected, actual)
	}
}

func TestWithSmartEmailStatus(t *testing.T) {
	const expected = transactional.ActiveSmartEmail
	ops := &transactional.Options{}
	option := transactional.WithSmartEmailStatus(expected)
	option(ops)
	actual := ops.SmartEmailStatus()
	if actual != expected {
		t.Errorf("Expected smart email status: %s, Actual: %s", expected, actual)
	}
}
