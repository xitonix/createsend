package createsend

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"testing"

	"github.com/xitonix/createsend/mock"
)

func checkRequestBody(t *testing.T, actual io.ReadCloser, expected *bodyMock) {
	t.Helper()
	defer actual.Close()
	b, err := ioutil.ReadAll(actual)
	if err != nil {
		t.Errorf("Did not expect an error but received: '%v'", err)
	}
	var body bodyMock
	err = json.Unmarshal(b, &body)
	if err != nil {
		t.Errorf("Did not expect an error but received: '%v'", err)
	}

	if expected == nil {
		if len(b) > 0 {
			t.Errorf("The request body should have been empty")
		}
		return
	}
	if expected.Value != body.Value {
		t.Errorf("Expected Request Body Value: %s, Actual: %s", expected.Value, body.Value)
	}
}

func checkErrorType(t *testing.T, err error, expectServerError bool) {
	t.Helper()
	var csErr *Error
	ok := errors.As(err, &csErr)
	if !ok {
		t.Error("We should always return a custom createsend Error type")
	}
	if csErr.IsFromServer() != expectServerError {
		t.Errorf("Expected server error: %v, actual: %v", expectServerError, csErr.IsFromServer())
	}
}

func createClient(t *testing.T, oAuthAuthentication, forceClientSideError bool) (*Client, *mock.HTTPClientMock) {
	t.Helper()
	httpClient := mock.NewHTTPClientMock(mock.ForceToFail(forceClientSideError))
	options := []Option{
		WithBaseURL("https://base.com"),
		WithHTTPClient(httpClient),
	}
	if oAuthAuthentication {
		options = append(options, WithOAuthToken("token"))
	} else {
		options = append(options, WithAPIKey("api_key"))
	}
	client, err := New(options...)
	if err != nil {
		checkErrorType(t, err, true)
		t.Fatalf("Did not expect an error but received: '%v'", err)
	}
	return client, httpClient
}

func checkError(actual, expected error) bool {
	if actual == nil {
		return expected == nil
	}
	if expected == nil {
		return false
	}
	return errors.Is(actual, expected)
}
