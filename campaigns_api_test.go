package createsend

import (
	"bytes"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/xitonix/createsend/campaigns"
	"github.com/xitonix/createsend/mock"
	"github.com/xitonix/createsend/order"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestCampaignsAPI_Create(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expectedError        error
		input                campaigns.WithURLs
		expectedResult       string
		oAuthAuthentication  bool
	}{
		{
			title: "successful execution",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`"campaign_id"`)),
			},
			expectedResult: "campaign_id",
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
			actual, err := client.Campaigns().Create("client_id", tC.input)
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceHTTPClientError)
			}
			if err == nil && tC.expectedError != nil {
				t.Errorf("Expected error '%v', received none", tC.expectedError)
			}
			if actual != tC.expectedResult {
				t.Errorf("Expected: %v, Actual: %v", tC.expectedResult, actual)
			}
		})
	}
}

func TestCampaignsAPI_Delete(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expectedError        error
		expectedResult       string
		oAuthAuthentication  bool
	}{
		{
			title: "successful execution",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`"campaign_id"`)),
			},
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
			httpClient.SetResponse("campaigns/campaign_id.json", tC.response)
			err := client.Campaigns().Delete("campaign_id")
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceHTTPClientError)
			}
			if err == nil && tC.expectedError != nil {
				t.Errorf("Expected error '%v', received none", tC.expectedError)
			}
		})
	}
}

func TestCampaignsAPI_Unschedule(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expectedError        error
		expectedResult       string
		oAuthAuthentication  bool
	}{
		{
			title: "successful execution",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`"campaign_id"`)),
			},
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
			httpClient.SetResponse("campaigns/campaign_id/unschedule.json", tC.response)
			err := client.Campaigns().Unschedule("campaign_id")
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceHTTPClientError)
			}
			if err == nil && tC.expectedError != nil {
				t.Errorf("Expected error '%v', received none", tC.expectedError)
			}
		})
	}
}

func TestCampaignsAPI_Send(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expectedError        error
		input                []string
		oAuthAuthentication  bool
	}{
		{
			title: "successful execution",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`"campaign_id"`)),
			},
			input: []string{"email1@testing.com", "email2@testing.com"},
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
			httpClient.SetResponse("campaigns/campaign_id/send.json", tC.response)
			err := client.Campaigns().Send("campaign_id", tC.input...)
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceHTTPClientError)
			}
			if err == nil && tC.expectedError != nil {
				t.Errorf("Expected error '%v', received none", tC.expectedError)
			}
		})
	}
}

func TestCampaignsAPI_SendAt(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expectedError        error
		input1               []string
		input2               time.Time
		oAuthAuthentication  bool
	}{
		{
			title: "successful execution",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`"campaign_id"`)),
			},
			input1: []string{"email1@testing.com", "email2@testing.com"},
			input2: time.Date(2060, 05, 20, 16, 45, 00, 00, time.UTC),
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
			httpClient.SetResponse("campaigns/campaign_id/send.json", tC.response)
			err := client.Campaigns().SendAt("campaign_id", tC.input2, tC.input1...)
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceHTTPClientError)
			}
			if err == nil && tC.expectedError != nil {
				t.Errorf("Expected error '%v', received none", tC.expectedError)
			}
		})
	}
}

func TestCampaignsAPI_SendPreview(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expectedError        error
		input                []string
		oAuthAuthentication  bool
	}{
		{
			title: "successful execution",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`"campaign_id"`)),
			},
			input: []string{"email1@testing.com", "email2@testing.com"},
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
			httpClient.SetResponse("campaigns/campaign_id/sendpreview.json", tC.response)
			err := client.Campaigns().SendPreview("campaign_id", tC.input...)
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceHTTPClientError)
			}
			if err == nil && tC.expectedError != nil {
				t.Errorf("Expected error '%v', received none", tC.expectedError)
			}
		})
	}
}

func TestCampaignsAPI_Summary(t *testing.T) {
	testCases := []struct {
		title                 string
		expectClientSideError bool
		response              *http.Response
		expected              campaigns.Summary
		expectedError         error
	}{
		{
			title: "no email client usage",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			},
			expected: campaigns.Summary{},
		},
		{
			title: "some email client usage",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"Recipients": 1000,
					"TotalOpened": 345,
					"Clicks": 132,
					"Unsubscribed": 43,
					"Bounced": 15,
					"UniqueOpened": 298,
					"SpamComplaints": 23,
					"WebVersionURL": "http://createsend.com/t/y-A1A1A1A1A1A1A1A1A1A1A1A1/",
					"WebVersionTextURL": "http://createsend.com/t/y-A1A1A1A1A1A1A1A1A1A1A1A1/t",
					"WorldviewURL": "http://myclient.createsend.com/reports/wv/y/8WY898U9U98U9U9",
					"Forwards": 18,
					"Likes": 25,
					"Mentions": 11
				}`)),
			},
			expected: campaigns.Summary{
				Recipients:        1000,
				TotalOpened:       345,
				Clicks:            132,
				Unsubscribed:      43,
				Bounced:           15,
				UniqueOpened:      298,
				SpamComplaints:    23,
				WebVersionURL:     "http://createsend.com/t/y-A1A1A1A1A1A1A1A1A1A1A1A1/",
				WebVersionTextURL: "http://createsend.com/t/y-A1A1A1A1A1A1A1A1A1A1A1A1/t",
				WorldviewURL:      "http://myclient.createsend.com/reports/wv/y/8WY898U9U98U9U9",
				Forwards:          18,
				Likes:             25,
				Mentions:          11,
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
			httpClient.SetResponse("campaigns/campaign_id/summary.json", tC.response)
			actual, err := client.Campaigns().Summary("campaign_id")
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

func TestCampaignsAPI_EmailClientUsage(t *testing.T) {
	testCases := []struct {
		title                 string
		expectClientSideError bool
		response              *http.Response
		expected              []campaigns.EmailClientUsage
		expectedError         error
	}{
		{
			title: "no email client usage",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`[]`)),
			},
			expected: []campaigns.EmailClientUsage{},
		},
		{
			title: "some email client usage",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
					{
						"Client": "iOS Devices",
						"Version": "iPhone",
						"Percentage": 19.83,
						"Subscribers": 7056
					},
					{
						"Client": "Apple Mail",
						"Version": "Apple Mail 6",
						"Percentage": 13.02,
						"Subscribers": 4633
					},
					{
						"Client": "Microsoft Outlook",
						"Version": "Outlook 2010",
						"Percentage": 7.18,
						"Subscribers": 2556
					},
					{
						"Client": "Undetectable",
						"Version": "Undetectable",
						"Percentage": 4.94,
						"Subscribers": 1632
					}
				]`)),
			},
			expected: []campaigns.EmailClientUsage{
				{
					Client:      "iOS Devices",
					Version:     "iPhone",
					Percentage:  19.83,
					Subscribers: 7056,
				},
				{
					Client:      "Apple Mail",
					Version:     "Apple Mail 6",
					Percentage:  13.02,
					Subscribers: 4633,
				},
				{
					Client:      "Microsoft Outlook",
					Version:     "Outlook 2010",
					Percentage:  7.18,
					Subscribers: 2556,
				},
				{
					Client:      "Undetectable",
					Version:     "Undetectable",
					Percentage:  4.94,
					Subscribers: 1632,
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
			httpClient.SetResponse("campaigns/campaign_id/emailclientusage.json", tC.response)
			actual, err := client.Campaigns().EmailClientUsage("campaign_id")
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

func TestCampaignsAPI_ListsAndSegments(t *testing.T) {
	testCases := []struct {
		title                 string
		expectClientSideError bool
		response              *http.Response
		expected              campaigns.ListsAndSegments
		expectedError         error
	}{
		{
			title: "no lists and segments for a given campaign",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"Lists": [],
					"Segments": []
				}`)),
			},
			expected: campaigns.ListsAndSegments{
				Lists:    []campaigns.List{},
				Segments: []campaigns.Segment{},
			},
		},
		{
			title: "some lists and segments for a given campaign",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"Lists": [
						{
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
							"Name": "My List 1"
						},
						{
							"ListID": "b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2",
							"Name": "My List 2"
						}
					],
					"Segments": [
						{
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
							"SegmentID": "c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3",
							"Title": "My Segment 1"
						},
						{
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
							"SegmentID": "d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4",
							"Title": "My Segment 2"
						}
					]
				}`)),
			},
			expected: campaigns.ListsAndSegments{
				Lists: []campaigns.List{
					{
						Name: "My List 1",
						ID:   "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
					},
					{
						Name: "My List 2",
						ID:   "b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2",
					},
				},
				Segments: []campaigns.Segment{
					{
						Title:  "My Segment 1",
						ListID: "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
						ID:     "c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3",
					},
					{
						Title:  "My Segment 2",
						ListID: "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
						ID:     "d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4",
					},
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
			httpClient.SetResponse("campaigns/campaign_id/listsandsegments.json", tC.response)
			actual, err := client.Campaigns().ListsAndSegments("campaign_id")
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

func TestCampaignsAPI_SentCampaignRecipients(t *testing.T) {
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

func TestCampaignsAPI_SentCampaignRecipients_DoesNotSupportDateOrder(t *testing.T) {
	expectedError := newClientError(ErrCodeInvalidDateOrderField)
	client, _ := createClient(t, true, false)
	_, err := client.Campaigns().Recipients("campaign_id", 1, 100, order.Date, order.DESC)
	if !checkError(err, expectedError) {
		t.Errorf("Expected '%v' error, actual: '%v'", expectedError, err)
	}
}

func TestCampaignsAPI_Bounces(t *testing.T) {
	testCases := []struct {
		title                 string
		expectClientSideError bool
		response              *http.Response
		expected              campaigns.Bounces
		expectedError         error
	}{
		{
			title: "no bounces",
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
			expected: campaigns.Bounces{
				Results:   []campaigns.Bounce{},
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
			title: "some bounces",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"Results": [
						{
							"EmailAddress": "example+1@example.com",
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
							"BounceType": "Hard",
							"Date": "2009-05-18 16:45:00",
							"Reason": "Invalid Email Address"
						},
						{
							"EmailAddress": "example+2@example.com",
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
							"BounceType": "Soft",
							"Date": "2009-05-20 16:45:00",
							"Reason": "Soft Bounce - Mailbox Full"
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
			expected: campaigns.Bounces{
				Results: []campaigns.Bounce{
					{
						Recipient: campaigns.Recipient{
							EmailAddress: "example+1@example.com",
							ListID:       "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
						},
						Date:       time.Date(2009, 05, 18, 16, 45, 00, 00, time.UTC),
						BounceType: "Hard",
						Reason:     "Invalid Email Address",
					},
					{
						Recipient: campaigns.Recipient{
							EmailAddress: "example+2@example.com",
							ListID:       "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
						},
						Date:       time.Date(2009, 05, 20, 16, 45, 00, 00, time.UTC),
						BounceType: "Soft",
						Reason:     "Soft Bounce - Mailbox Full",
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
			httpClient.SetResponse("campaigns/campaign_id/bounces.json", tC.response)
			bouncesSinceDate := time.Date(2000, 05, 20, 16, 45, 00, 00, time.UTC)
			actual, err := client.Campaigns().Bounces("campaign_id", bouncesSinceDate, 1, 100, order.Date, order.DESC)
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

func TestCampaignsAPI_BouncesDateParseError(t *testing.T) {
	testCases := []struct {
		title         string
		response      *http.Response
		expectedError error
	}{
		{
			title: "invalid date to parse",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"Results": [
						{
							"EmailAddress": "example+1@example.com",
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
							"BounceType": "Hard",
							"Date": "XYZ",
							"Reason": "Invalid Email Address"
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
			expectedError: fmt.Errorf("Could not find format for %q", "XYZ"),
		},
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {
			client, httpClient := createClient(t, true, false)
			httpClient.SetResponse("campaigns/campaign_id/bounces.json", tC.response)
			_, err := client.Campaigns().Bounces("campaign_id", time.Time{}, 1, 100, order.Date, order.DESC)
			if err == nil {
				t.Error("Expected date parse error, got none")
				return
			}
			if err.Error() != tC.expectedError.Error() {
				t.Errorf("Expected '%v', actual '%v'", tC.expectedError, err)
			}
		})
	}
}

func TestCampaignsAPI_Unsubscribes(t *testing.T) {
	testCases := []struct {
		title                 string
		expectClientSideError bool
		response              *http.Response
		expected              campaigns.Unsubscribes
		expectedError         error
	}{
		{
			title: "no unsubscribes",
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
			expected: campaigns.Unsubscribes{
				Results:   []campaigns.Unsubscribe{},
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
			title: "some unsubscribes",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"Results": [
						{
							"EmailAddress": "example+1@example.com",
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
							"Date": "2009-05-18 16:45:00",
							"IPAddress": "192.168.0.1"
						},
						{
							"EmailAddress": "example+1@example.com",
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
							"Date": "2009-05-19 16:45:00",
							"IPAddress": "192.168.0.1"
						},
						{
							"EmailAddress": "example+2@example.com",
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
							"Date": "2009-05-20 16:45:00",
							"IPAddress": "192.168.0.3"
						},
						{
							"EmailAddress": "example+3@example.com",
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
							"Date": "2009-05-21 16:45:00",
							"IPAddress": "192.168.0.4"
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
			expected: campaigns.Unsubscribes{
				Results: []campaigns.Unsubscribe{
					{
						Recipient: campaigns.Recipient{
							EmailAddress: "example+1@example.com",
							ListID:       "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
						},
						Date:      time.Date(2009, 05, 18, 16, 45, 00, 00, time.UTC),
						IPAddress: "192.168.0.1",
					},
					{
						Recipient: campaigns.Recipient{
							EmailAddress: "example+1@example.com",
							ListID:       "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
						},
						Date:      time.Date(2009, 05, 19, 16, 45, 00, 00, time.UTC),
						IPAddress: "192.168.0.1",
					},
					{
						Recipient: campaigns.Recipient{
							EmailAddress: "example+2@example.com",
							ListID:       "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
						},
						Date:      time.Date(2009, 05, 20, 16, 45, 00, 00, time.UTC),
						IPAddress: "192.168.0.3",
					},
					{
						Recipient: campaigns.Recipient{
							EmailAddress: "example+3@example.com",
							ListID:       "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
						},
						Date:      time.Date(2009, 05, 21, 16, 45, 00, 00, time.UTC),
						IPAddress: "192.168.0.4",
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
			httpClient.SetResponse("campaigns/campaign_id/unsubscribes.json", tC.response)
			actual, err := client.Campaigns().Unsubscribes("campaign_id", time.Time{}, 1, 100, order.Date, order.DESC)
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

func TestCampaignsAPI_UnsubscribesDateParseError(t *testing.T) {
	testCases := []struct {
		title         string
		response      *http.Response
		expectedError error
	}{
		{
			title: "invalid date to parse",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"Results": [
						{
							"EmailAddress": "example+1@example.com",
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
							"Date": "XYZ",
							"IPAddress": "192.168.0.1"
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
			expectedError: fmt.Errorf("Could not find format for %q", "XYZ"),
		},
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {
			client, httpClient := createClient(t, true, false)
			httpClient.SetResponse("campaigns/campaign_id/unsubscribes.json", tC.response)
			_, err := client.Campaigns().Unsubscribes("campaign_id", time.Time{}, 1, 100, order.Date, order.DESC)
			if err == nil {
				t.Error("Expected date parse error, got none")
				return
			}
			if err.Error() != tC.expectedError.Error() {
				t.Errorf("Expected '%v', actual '%v'", tC.expectedError, err)
			}
		})
	}
}

func TestCampaignsAPI_Opens(t *testing.T) {
	testCases := []struct {
		title                 string
		expectClientSideError bool
		response              *http.Response
		expected              campaigns.Opens
		expectedError         error
	}{
		{
			title: "no opens",
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
			expected: campaigns.Opens{
				Results:   []campaigns.OpenDetails{},
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
			title: "some opens",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"Results": [
						{
							"EmailAddress": "example+1@example.com",
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
							"Date": "2009-05-18 16:45:00",
							"IPAddress": "192.168.0.1",
							"Latitude": -33.8683,
							"Longitude": 151.2086,
							"City": "Sydney",
							"Region": "New South Wales",
							"CountryCode": "AU",
							"CountryName": "Australia"
						},
						{
							"EmailAddress": "example+2@example.com",
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
							"Date": "2009-05-20 16:45:00",
							"IPAddress": "192.168.0.3",
							"Latitude": -33.8683,
							"Longitude": 151.2086,
							"City": "Sydney",
							"Region": "New South Wales",
							"CountryCode": "AU",
							"CountryName": "Australia"
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
			expected: campaigns.Opens{
				Results: []campaigns.OpenDetails{
					{
						Recipient: campaigns.Recipient{
							EmailAddress: "example+1@example.com",
							ListID:       "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
						},
						RecipientLocationDetails: campaigns.RecipientLocationDetails{
							IPAddress:   "192.168.0.1",
							Latitude:    -33.8683,
							Longitude:   151.2086,
							City:        "Sydney",
							Region:      "New South Wales",
							CountryCode: "AU",
							CountryName: "Australia",
						},
						Date: time.Date(2009, 05, 18, 16, 45, 00, 00, time.UTC),
					},
					{
						Recipient: campaigns.Recipient{
							EmailAddress: "example+2@example.com",
							ListID:       "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
						},
						RecipientLocationDetails: campaigns.RecipientLocationDetails{
							IPAddress:   "192.168.0.3",
							Latitude:    -33.8683,
							Longitude:   151.2086,
							City:        "Sydney",
							Region:      "New South Wales",
							CountryCode: "AU",
							CountryName: "Australia",
						},
						Date: time.Date(2009, 05, 20, 16, 45, 00, 00, time.UTC),
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
			httpClient.SetResponse("campaigns/campaign_id/opens.json", tC.response)
			actual, err := client.Campaigns().Opens("campaign_id", time.Time{}, 1, 100, order.Date, order.DESC)
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

func TestCampaignsAPI_OpensDateParseError(t *testing.T) {
	testCases := []struct {
		title         string
		response      *http.Response
		expectedError error
	}{
		{
			title: "invalid date to parse",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"Results": [
						{
							"EmailAddress": "example+1@example.com",
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
							"Date": "XYZ",
							"IPAddress": "192.168.0.1",
							"Latitude": -33.8683,
							"Longitude": 151.2086,
							"City": "Sydney",
							"Region": "New South Wales",
							"CountryCode": "AU",
							"CountryName": "Australia"
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
			expectedError: fmt.Errorf("Could not find format for %q", "XYZ"),
		},
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {
			client, httpClient := createClient(t, true, false)
			httpClient.SetResponse("campaigns/campaign_id/opens.json", tC.response)
			_, err := client.Campaigns().Opens("campaign_id", time.Time{}, 1, 100, order.Date, order.DESC)
			if err == nil {
				t.Error("Expected date parse error, got none")
				return
			}
			if err.Error() != tC.expectedError.Error() {
				t.Errorf("Expected '%v', actual '%v'", tC.expectedError, err)
			}
		})
	}
}

func TestCampaignsAPI_Clicks(t *testing.T) {
	testCases := []struct {
		title                 string
		expectClientSideError bool
		response              *http.Response
		expected              campaigns.Clicks
		expectedError         error
	}{
		{
			title: "no clicks",
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
			expected: campaigns.Clicks{
				Results:   []campaigns.ClickDetails{},
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
			title: "some clicks",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"Results": [
						{
							"EmailAddress": "example+1@example.com",
							"URL": "http://www.myexammple.com/index.html",
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
							"Date": "2009-05-18 16:45:00",
							"IPAddress": "192.168.0.1",
							"Latitude": -33.8683,
							"Longitude": 151.2086,
							"City": "Sydney",
							"Region": "New South Wales",
							"CountryCode": "AU",
							"CountryName": "Australia"
						},
						{
							"EmailAddress": "example+2@example.com",
							"URL": "http://www.myexammple.com/index.html",
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
							"Date": "2009-05-20 16:45:00",
							"IPAddress": "192.168.0.3",
							"Latitude": -33.8683,
							"Longitude": 151.2086,
							"City": "Sydney",
							"Region": "New South Wales",
							"CountryCode": "AU",
							"CountryName": "Australia"
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
			expected: campaigns.Clicks{
				Results: []campaigns.ClickDetails{
					{
						Recipient: campaigns.Recipient{
							EmailAddress: "example+1@example.com",
							ListID:       "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
						},
						RecipientLocationDetails: campaigns.RecipientLocationDetails{
							IPAddress:   "192.168.0.1",
							Latitude:    -33.8683,
							Longitude:   151.2086,
							City:        "Sydney",
							Region:      "New South Wales",
							CountryCode: "AU",
							CountryName: "Australia",
						},
						URL:  "http://www.myexammple.com/index.html",
						Date: time.Date(2009, 05, 18, 16, 45, 00, 00, time.UTC),
					},
					{
						Recipient: campaigns.Recipient{
							EmailAddress: "example+2@example.com",
							ListID:       "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
						},
						RecipientLocationDetails: campaigns.RecipientLocationDetails{
							IPAddress:   "192.168.0.3",
							Latitude:    -33.8683,
							Longitude:   151.2086,
							City:        "Sydney",
							Region:      "New South Wales",
							CountryCode: "AU",
							CountryName: "Australia",
						},
						URL:  "http://www.myexammple.com/index.html",
						Date: time.Date(2009, 05, 20, 16, 45, 00, 00, time.UTC),
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
			httpClient.SetResponse("campaigns/campaign_id/clicks.json", tC.response)
			actual, err := client.Campaigns().Clicks("campaign_id", time.Time{}, 1, 100, order.Date, order.DESC)
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

func TestCampaignsAPI_ClicksDateParseError(t *testing.T) {
	testCases := []struct {
		title         string
		response      *http.Response
		expectedError error
	}{
		{
			title: "invalid date to parse",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"Results": [
						{
							"EmailAddress": "example+1@example.com",
							"URL": "http://www.myexammple.com/index.html",
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
							"Date": "XYZ",
							"IPAddress": "192.168.0.1",
							"Latitude": -33.8683,
							"Longitude": 151.2086,
							"City": "Sydney",
							"Region": "New South Wales",
							"CountryCode": "AU",
							"CountryName": "Australia"
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
			expectedError: fmt.Errorf("Could not find format for %q", "XYZ"),
		},
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {
			client, httpClient := createClient(t, true, false)
			httpClient.SetResponse("campaigns/campaign_id/clicks.json", tC.response)
			_, err := client.Campaigns().Clicks("campaign_id", time.Time{}, 1, 100, order.Date, order.DESC)
			if err == nil {
				t.Error("Expected date parse error, got none")
				return
			}
			if err.Error() != tC.expectedError.Error() {
				t.Errorf("Expected '%v', actual '%v'", tC.expectedError, err)
			}
		})
	}
}

func TestCampaignsAPI_SpamComplaints(t *testing.T) {
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

func TestCampaignsAPI_SpamComplaintsDateParseError(t *testing.T) {
	testCases := []struct {
		title         string
		response      *http.Response
		expectedError error
	}{
		{
			title: "invalid date to parse",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"Results": [
						{
							"EmailAddress": "example+1@example.com",
							"ListID": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
							"Date": "XYZ"
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
			expectedError: fmt.Errorf("Could not find format for %q", "XYZ"),
		},
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {
			client, httpClient := createClient(t, true, false)
			httpClient.SetResponse("campaigns/campaign_id/spam.json", tC.response)
			_, err := client.Campaigns().SpamComplaints("campaign_id", time.Time{}, 1, 100, order.Date, order.DESC)
			if err == nil {
				t.Error("Expected date parse error, got none")
				return
			}
			if err.Error() != tC.expectedError.Error() {
				t.Errorf("Expected '%v', actual '%v'", tC.expectedError, err)
			}
		})
	}
}
