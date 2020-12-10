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
				if _, ok := err.(*createsend.Error); !ok {
					t.Errorf("We should have returned a custom createsend error")
				}
				return
			}

			if client.Accounts() == nil {
				t.Errorf("Accounts API should not be nil")
			}

			if client.Clients() == nil {
				t.Errorf("Clients API should not be nil")
			}

			if client.Campaigns() == nil {
				t.Errorf("Campaigns API should not be nil")
			}
		})
	}
}
