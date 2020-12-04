package createsend

import (
	"context"
	"net/http"
	"testing"
)

func TestDefaultOptions(t *testing.T) {
	ops := defaultOptions()
	if ops.baseURL != DefaultBaseURL {
		t.Errorf("Expected base URL: %s, Actual: %s", DefaultBaseURL, ops.baseURL)
	}

	if ops.auth.method != undefinedAuthentication {
		t.Errorf("Expected authentication method: %v, Actual: %v", undefinedAuthentication, ops.auth.method)
	}

	if ops.auth.token != "" {
		t.Errorf("Expected the default authentication token to be empty, Actual: %s", ops.auth.token)
	}

	if ops.client == nil {
		t.Error("The default HTTP client was nil")
	}

	if ops.ctx == nil {
		t.Error("The default context was nil")
	}
}

func TestWithContext(t *testing.T) {
	testCases := []struct {
		title string
		ctx   context.Context
	}{
		{
			title: "setting the context to nil must use the background context",
		},
		{
			title: "non empty context",
			ctx:   context.TODO(),
		},
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {
			ops := defaultOptions()
			option := WithContext(tC.ctx)
			option(ops)
			if ops.ctx == nil {
				t.Error("The context was nil")
			}
		})
	}
}

func TestWithClientsAPI(t *testing.T) {
	ops := defaultOptions()
	option := WithClientsAPI(&clientsAPI{})
	option(ops)
	if ops.clients == nil {
		t.Error("Clients API was nil")
	}
}

func TestWithAccountsAPI(t *testing.T) {
	ops := defaultOptions()
	option := WithAccountsAPI(&accountsAPI{})
	option(ops)
	if ops.accounts == nil {
		t.Error("Accounts API was nil")
	}
}

func TestWithTransactionalAPI(t *testing.T) {
	ops := defaultOptions()
	option := WithTransactionalAPI(&transactionalAPI{})
	option(ops)
	if ops.transactional == nil {
		t.Error("Transactional API was nil")
	}
}

func TestWithHTTPClient(t *testing.T) {
	ops := defaultOptions()
	option := WithHTTPClient(&http.Client{})
	option(ops)
	if ops.client == nil {
		t.Error("HTTP client was nil")
	}
}

func TestWithBaseURL(t *testing.T) {
	expected := "base"
	ops := defaultOptions()
	option := WithBaseURL(expected)
	option(ops)
	if ops.baseURL != expected {
		t.Errorf("Expected base URL: %s, Actual: %s", expected, ops.baseURL)
	}
}

func TestWithAPIKey(t *testing.T) {
	expected := "key"
	ops := defaultOptions()
	option := WithAPIKey(expected)
	option(ops)
	if ops.auth.method != apiKeyAuthentication {
		t.Errorf("Expected authentication method: %v, Actual: %v", apiKeyAuthentication, ops.auth.method)
	}

	if ops.auth.token != expected {
		t.Errorf("Expected authentication token: %s, Actual: %s", expected, ops.auth.token)
	}
}

func TestWithOAuthToken(t *testing.T) {
	expected := "token"
	ops := defaultOptions()
	option := WithOAuthToken(expected)
	option(ops)
	if ops.auth.method != oAuthAuthentication {
		t.Errorf("Expected authentication method: %v, Actual: %v", oAuthAuthentication, ops.auth.method)
	}

	if ops.auth.token != expected {
		t.Errorf("Expected authentication token: %s, Actual: %s", expected, ops.auth.token)
	}
}
