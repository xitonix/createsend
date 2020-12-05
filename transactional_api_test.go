package createsend

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/xitonix/createsend/mock"
	"github.com/xitonix/createsend/transactional"
)

func TestTransactionalAPI_SmartEmails(t *testing.T) {
	date := time.Date(2020, 12, 1, 20, 21, 22, 0, time.UTC)
	testCases := []struct {
		title                 string
		clientID              string
		smartEmailStatus      transactional.SmartEmailStatus
		forceHTTPClientError  bool
		expectClientSideError bool
		response              *http.Response
		expected              []*transactional.SmartEmailBasicDetails
		expectedError         error
		oAuthAuthentication   bool
	}{
		{
			title: "no smart emails",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`[]`)),
			},
			expected: []*transactional.SmartEmailBasicDetails{},
		},
		{
			title: "no smart emails and empty server response body",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: []*transactional.SmartEmailBasicDetails{},
		},
		{
			title: "with smart emails",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				 {
					"ID": "id1",
					"Name": "name1",
					"CreatedAt": "2020-12-01T20:21:22",
					"Status": "Active"
				},
				{
					"ID": "id2",
					"Name": "name2",
					"CreatedAt": "2020-12-01T20:21:22",
					"Status": "Draft"
				}
			]`)),
			},
			expected: []*transactional.SmartEmailBasicDetails{
				{
					ID:        "id1",
					Name:      "name1",
					CreatedAt: date,
					Status:    transactional.ActiveSmartEmail,
				},
				{
					ID:        "id2",
					Name:      "name2",
					CreatedAt: date,
					Status:    transactional.DraftSmartEmail,
				},
			},
		},
		{
			title: "invalid date value",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"ID": "id",
					"Name": "name",
					"CreatedAt": "invalid date",
					"Status": "Active"
				}
			]`)),
			},
			expected:              nil,
			expectedError:         newClientError(ErrCodeDataProcessing),
			expectClientSideError: true,
		},
		{
			title: "unknown smart email status",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"ID": "id",
					"Name": "name",
					"CreatedAt": "2020-12-01T20:21:22",
					"Status": "what?"
				}
			]`)),
			},
			expected: []*transactional.SmartEmailBasicDetails{
				{
					ID:        "id",
					Name:      "name",
					CreatedAt: date,
					Status:    transactional.UnknownSmartEmail,
				},
			},
		},
		{
			title:            "filtered by smart email status",
			smartEmailStatus: transactional.DraftSmartEmail,
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"ID": "id",
					"Name": "name",
					"CreatedAt": "2020-12-01T20:21:22",
					"Status": "draft"
				}
			]`)),
			},
			expected: []*transactional.SmartEmailBasicDetails{
				{
					ID:        "id",
					Name:      "name",
					CreatedAt: date,
					Status:    transactional.DraftSmartEmail,
				},
			},
		},
		{
			title:    "query with client ID",
			clientID: "client_id",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"ID": "id",
					"Name": "name",
					"CreatedAt": "2020-12-01T20:21:22",
					"Status": "active"
				}
			]`)),
			},
			expected: []*transactional.SmartEmailBasicDetails{
				{
					ID:        "id",
					Name:      "name",
					CreatedAt: date,
					Status:    transactional.ActiveSmartEmail,
				},
			},
		},
		{
			title: "oAuth Authentication",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"ID": "id",
					"Name": "name",
					"CreatedAt": "2020-12-01T20:21:22",
					"Status": "Active"
				}
			]`)),
			},
			expected: []*transactional.SmartEmailBasicDetails{
				{
					ID:        "id",
					Name:      "name",
					CreatedAt: date,
					Status:    transactional.ActiveSmartEmail,
				},
			},
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
			httpClient.SetResponse("transactional/smartEmail", tC.response)
			actual, err := client.Transactional().SmartEmails(transactional.WithClientID(tC.clientID), transactional.WithSmartEmailStatus(tC.smartEmailStatus))
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.expectClientSideError && !tC.forceHTTPClientError)
			}

			expectedStatus := "all"
			if tC.smartEmailStatus != transactional.UnknownSmartEmail {
				expectedStatus = tC.smartEmailStatus.String()
			}
			expectedQuery := map[string]string{
				"status": expectedStatus,
			}
			if tC.clientID != "" {
				expectedQuery["clientID"] = tC.clientID
			}

			checkQueryStringParameters(t, httpClient.LastRequest(), expectedQuery)

			if diff := cmp.Diff(tC.expected, actual); diff != "" {
				t.Errorf("Expectations failed (-expected +actual):\n%s", diff)
			}
		})
	}
}
