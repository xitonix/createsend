package createsend_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/xitonix/createsend"
	"github.com/xitonix/createsend/accounts"
	"github.com/xitonix/createsend/internal/test"
	"github.com/xitonix/createsend/mock"
)

func TestClients(t *testing.T) {
	testCases := []struct {
		title                 string
		simulateServerFailure bool
		response              *http.Response
		expected              []*accounts.Client
		expectedError         error
	}{
		{
			title: "account with no clients",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`[]`)),
			},
			expected: []*accounts.Client{},
		},
		{
			title: "account with no clients and empty server response body",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: []*accounts.Client{},
		},
		{
			title: "account with clients",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`[{"ClientID":"id", "Name":"name"}]`)),
			},
			expected: []*accounts.Client{{
				ClientID: "id",
				Name:     "name",
			}},
		},
		{
			title:                 "simulate remote call failure",
			response:              &http.Response{},
			simulateServerFailure: true,
			expectedError:         mock.ErrDeliberate,
		},
		{
			title: "simulate server side error",
			response: &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":500}`)),
			},
			expectedError: &createsend.Error{Code: 500},
		},
	}

	path := "/clients.json"

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {
			httpClient := mock.NewHTTPClientMock(mock.ForceToFail(tC.simulateServerFailure))
			client, err := createsend.New(
				createsend.WithBaseURL("https://base.com"),
				createsend.WithHTTPClient(httpClient),
				createsend.WithAPIKey("api_key"),
			)
			if err != nil {
				t.Errorf("Did not expect an error but received: '%v'", err)
				checkErrorType(t, err)
			}
			httpClient.SetResponse(path, tC.response)
			actual, err := client.Accounts().Clients()
			if err != nil {
				if !test.CheckError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err)
			}
			if !reflect.DeepEqual(actual, tC.expected) {
				t.Errorf("Expected '%+v', actual: '%+v'", tC.expected, actual)
			}
		})
	}
}

func TestBillingDetails(t *testing.T) {
	testCases := []struct {
		title                 string
		simulateServerFailure bool
		response              *http.Response
		expected              *accounts.Billing
		expectedError         error
	}{
		{
			title: "account with no billing details",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			},
			expected: &accounts.Billing{},
		},
		{
			title: "account with no billing details and empty server response body",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: nil,
		},
		{
			title: "account with billing details",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Credits":100}`)),
			},
			expected: &accounts.Billing{
				Credits: 100,
			},
		},
		{
			title:                 "simulate remote call failure",
			response:              &http.Response{},
			simulateServerFailure: true,
			expectedError:         mock.ErrDeliberate,
		},
		{
			title: "simulate server side error",
			response: &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":500}`)),
			},
			expectedError: &createsend.Error{Code: 500},
		},
	}

	path := "/billingdetails.json"

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {
			httpClient := mock.NewHTTPClientMock(mock.ForceToFail(tC.simulateServerFailure))
			client, err := createsend.New(
				createsend.WithBaseURL("https://base.com"),
				createsend.WithHTTPClient(httpClient),
				createsend.WithAPIKey("api_key"),
			)
			if err != nil {
				t.Errorf("Did not expect an error but received: '%v'", err)
				checkErrorType(t, err)
			}
			httpClient.SetResponse(path, tC.response)
			actual, err := client.Accounts().Billing()
			if err != nil {
				if !test.CheckError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err)
			}
			if !reflect.DeepEqual(actual, tC.expected) {
				t.Errorf("Expected '%+v', actual: '%+v'", tC.expected, actual)
			}
		})
	}
}
