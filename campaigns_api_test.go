package createsend

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"github.com/xitonix/createsend/campaigns"
	"github.com/xitonix/createsend/order"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestClientsAPI_SentCampaignRecipients(t *testing.T) {
	testCases := []struct {
		title                 string
		expectClientSideError bool
		response              *http.Response
		expected              campaigns.Recipients
		expectedError         error
	}{
		{
			title: "no sent campaign recipients",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"Results": [],
					"ResultsOrderedBy": "list",
					"OrderDirection": "desc",
					"PageNumber": 1,
					"PageSize": 100,
					"RecordsOnThisPage": 0,
					"TotalNumberOfRecords": 0,
					"NumberOfPages": 0
				}`)),
			},
			expected: campaigns.Recipients{
				Results:   []campaigns.Recipient{},
				OrderedBy: order.List,
				Page: order.Page{
					OrderDirection: order.DESC,
					Number:         1,
					Size:           100,
					Records:        0,
					Total:          0,
					NumberOfPages:  0,
				},
			},
		},
		{
			title: "some sent campaign recipients",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"Results": [
						{
							"EmailAddress": "example+1@example.com",
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1"
						},
						{
							"EmailAddress": "example+1@example.com",
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1"
						},
						{
							"EmailAddress": "example+2@example.com",
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1"
						},
						{
							"EmailAddress": "example+3@example.com",
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1"
						}
					],
					"ResultsOrderedBy": "list",
					"OrderDirection": "desc",
					"PageNumber": 1,
					"PageSize": 100,
					"RecordsOnThisPage": 4,
					"TotalNumberOfRecords": 4,
					"NumberOfPages": 1
				}`)),
			},
			expected: campaigns.Recipients{
				Results: []campaigns.Recipient{
					{
						EmailAddress: "example+1@example.com",
						ListID:       "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
					},
					{
						EmailAddress: "example+1@example.com",
						ListID:       "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
					},
					{
						EmailAddress: "example+2@example.com",
						ListID:       "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
					},
					{
						EmailAddress: "example+3@example.com",
						ListID:       "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
					},
				},
				OrderedBy: order.List,
				Page: order.Page{
					OrderDirection: order.DESC,
					Number:         1,
					Size:           100,
					Records:        4,
					Total:          4,
					NumberOfPages:  1,
				},
			},
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
			client, httpClient := createClient(t, true, false)
			httpClient.SetResponse("campaigns/campaign_id/recipients.json", tC.response)
			actual, err := client.Campaigns().Recipients("campaign_id", 1, 100, order.List, order.DESC)
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, true)
			}
			if diff := cmp.Diff(tC.expected, actual); diff != "" {
				t.Errorf("Expectations failed (-expected +actual):\n%s", diff)
			}
		})
	}
}

func TestClientsAPI_SpamComplaints(t *testing.T) {
	testCases := []struct {
		title                 string
		expectClientSideError bool
		response              *http.Response
		expected              campaigns.SpamComplaints
		expectedError         error
	}{
		{
			title: "no spam complaints",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"Results": [],
					"ResultsOrderedBy": "date",
					"OrderDirection": "desc",
					"PageNumber": 1,
					"PageSize": 100,
					"RecordsOnThisPage": 0,
					"TotalNumberOfRecords": 0,
					"NumberOfPages": 0
				}`)),
			},
			expected: campaigns.SpamComplaints{
				Results:   []campaigns.SpamComplaint{},
				OrderedBy: order.Date,
				Page: order.Page{
					OrderDirection: order.DESC,
					Number:         1,
					Size:           100,
					Records:        0,
					Total:          0,
					NumberOfPages:  0,
				},
			},
		},
		{
			title: "some spam complaints",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"Results": [
						{
							"EmailAddress": "example+1@example.com",
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
							"Date": "2009-05-18 16:45:00"
						},
						{
							"EmailAddress": "example+1@example.com",
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
							"Date": "2009-05-19 16:45:00"
						},
						{
							"EmailAddress": "example+2@example.com",
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
							"Date": "2009-05-20 16:45:00"
						},
						{
							"EmailAddress": "example+3@example.com",
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
							"Date": "2009-05-21 16:45:00"
						}
					],
					"ResultsOrderedBy": "date",
					"OrderDirection": "desc",
					"PageNumber": 1,
					"PageSize": 100,
					"RecordsOnThisPage": 4,
					"TotalNumberOfRecords": 4,
					"NumberOfPages": 1
				}`)),
			},
			expected: campaigns.SpamComplaints{
				Results: []campaigns.SpamComplaint{
					{
						Recipient: campaigns.Recipient{
							EmailAddress: "example+1@example.com",
							ListID:       "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
						},
						Date: time.Date(2009, 05, 18, 16, 45, 00, 00, time.UTC),
					},
					{
						Recipient: campaigns.Recipient{
							EmailAddress: "example+1@example.com",
							ListID:       "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
						},
						Date: time.Date(2009, 05, 19, 16, 45, 00, 00, time.UTC),
					},
					{
						Recipient: campaigns.Recipient{
							EmailAddress: "example+2@example.com",
							ListID:       "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
						},
						Date: time.Date(2009, 05, 20, 16, 45, 00, 00, time.UTC),
					},
					{
						Recipient: campaigns.Recipient{
							EmailAddress: "example+3@example.com",
							ListID:       "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
						},
						Date: time.Date(2009, 05, 21, 16, 45, 00, 00, time.UTC),
					},
				},
				OrderedBy: order.Date,
				Page: order.Page{
					OrderDirection: order.DESC,
					Number:         1,
					Size:           100,
					Records:        4,
					Total:          4,
					NumberOfPages:  1,
				},
			},
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
			client, httpClient := createClient(t, true, false)
			httpClient.SetResponse("campaigns/campaign_id/spam.json", tC.response)
			actual, err := client.Campaigns().SpamComplaints("campaign_id", time.Time{}, 1, 100, order.Date, order.DESC)
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, true)
			}
			if diff := cmp.Diff(tC.expected, actual); diff != "" {
				t.Errorf("Expectations failed (-expected +actual):\n%s", diff)
			}
		})
	}
}
