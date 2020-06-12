package createsend_test

import (
	"testing"

	"github.com/xitonix/createsend"
	"github.com/xitonix/createsend/mock"
)

func TestNewClient(t *testing.T) {
	testCases := []struct {
		title         string
		httpClient    createsend.HTTPClient
		expectedError error
	}{
		{
			title:      "successful initialisation",
			httpClient: mock.NewHTTPClientMock(),
		},
		{
			title:         "force internal client creation failure",
			httpClient:    nil,
			expectedError: &createsend.Error{},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {
			client, err := createsend.New(
				createsend.WithBaseURL("https://base.com"),
				createsend.WithHTTPClient(tC.httpClient),
				createsend.WithAPIKey("api_key"),
			)
			if err != nil {
				checkErrorType(t, err)
				return
			}

			if client.Accounts() == nil {
				t.Errorf("Account API should not be nil")
			}
		})
	}
}
