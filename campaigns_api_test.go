package createsend

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/xitonix/createsend/campaigns"
	"github.com/xitonix/createsend/mock"
)

func TestCampaignsAPI_CreateDraft(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expectedError        error
		oAuthAuthentication  bool
		expected             string
	}{
		{
			title: "empty JSON response",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
			},
			expected: "",
		},
		{
			title: "non empty server response",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`"campaign_id"`)),
			},
			expected: "campaign_id",
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`"campaign_id"`)),
			},
			expected:            "campaign_id",
			oAuthAuthentication: true,
		},
		{
			title:                "simulate remote call failure",
			response:             &http.Response{},
			forceHTTPClientError: true,
			expectedError:        mock.ErrDeliberate,
		},
		{
			title: "simulate server side error",
			response: &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":500}`)),
			},
			expectedError: &Error{Code: 500},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {
			client, httpClient := createClient(t, tC.oAuthAuthentication, tC.forceHTTPClientError)
			httpClient.SetResponse("campaigns/client_id.json", tC.response)
			actual, err := client.Campaigns().CreateDraft("client_id", campaigns.Draft{})
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceHTTPClientError)
			}
			if diff := cmp.Diff(tC.expected, actual); diff != "" {
				t.Errorf("Expectations failed (-expected +actual):\n%s", diff)
			}
		})
	}
}
