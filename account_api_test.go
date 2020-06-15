package createsend

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/xitonix/createsend/accounts"
	"github.com/xitonix/createsend/internal/test"
	"github.com/xitonix/createsend/mock"
)

func TestAccountsAPI_Clients(t *testing.T) {
	testCases := []struct {
		title                string
		forceClientSideError bool
		response             *http.Response
		expected             []*accounts.Client
		expectedError        error
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
			title:                "simulate remote call failure",
			response:             &http.Response{},
			forceClientSideError: true,
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
			httpClient := mock.NewHTTPClientMock(mock.ForceToFail(tC.forceClientSideError))
			client, err := New(
				WithBaseURL("https://base.com"),
				WithHTTPClient(httpClient),
				WithAPIKey("api_key"),
			)
			if err != nil {
				t.Errorf("Did not expect an error but received: '%v'", err)
				checkErrorType(t, err, true)
			}
			httpClient.SetResponse(listClientsPath, tC.response)
			actual, err := client.Accounts().Clients()
			if err != nil {
				if !test.CheckError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceClientSideError)
			}
			if !reflect.DeepEqual(actual, tC.expected) {
				t.Errorf("Expected '%+v', actual: '%+v'", tC.expected, actual)
			}
		})
	}
}

func TestAccountsAPI_Billing(t *testing.T) {
	testCases := []struct {
		title                string
		forceClientSideError bool
		response             *http.Response
		expected             *accounts.Billing
		expectedError        error
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
			title:                "simulate remote call failure",
			response:             &http.Response{},
			forceClientSideError: true,
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
			httpClient := mock.NewHTTPClientMock(mock.ForceToFail(tC.forceClientSideError))
			client, err := New(
				WithBaseURL("https://base.com"),
				WithHTTPClient(httpClient),
				WithAPIKey("api_key"),
			)
			if err != nil {
				t.Errorf("Did not expect an error but received: '%v'", err)
				checkErrorType(t, err, true)
			}
			httpClient.SetResponse(fetchBillingDetailsPath, tC.response)
			actual, err := client.Accounts().Billing()
			if err != nil {
				if !test.CheckError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceClientSideError)
			}
			if !reflect.DeepEqual(actual, tC.expected) {
				t.Errorf("Expected '%+v', actual: '%+v'", tC.expected, actual)
			}
		})
	}
}

func TestAccountsAPI_Countries(t *testing.T) {
	testCases := []struct {
		title                string
		forceClientSideError bool
		response             *http.Response
		expected             []string
		expectedError        error
	}{
		{
			title: "account with no valid countries",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`[]`)),
			},
			expected: []string{},
		},
		{
			title: "account with no valid countries and empty server response body",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: nil,
		},
		{
			title: "account with some valid countries",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`["country"]`)),
			},
			expected: []string{"country"},
		},
		{
			title:                "simulate remote call failure",
			response:             &http.Response{},
			forceClientSideError: true,
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
			httpClient := mock.NewHTTPClientMock(mock.ForceToFail(tC.forceClientSideError))
			client, err := New(
				WithBaseURL("https://base.com"),
				WithHTTPClient(httpClient),
				WithAPIKey("api_key"),
			)
			if err != nil {
				t.Errorf("Did not expect an error but received: '%v'", err)
				checkErrorType(t, err, true)
			}
			httpClient.SetResponse(fetchValidCountriesPath, tC.response)
			actual, err := client.Accounts().Countries()
			if err != nil {
				if !test.CheckError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceClientSideError)
			}
			if !reflect.DeepEqual(actual, tC.expected) {
				t.Errorf("Expected '%+v', actual: '%+v'", tC.expected, actual)
			}
		})
	}
}

func TestAccountsAPI_Timezones(t *testing.T) {
	testCases := []struct {
		title                string
		forceClientSideError bool
		response             *http.Response
		expected             []string
		expectedError        error
	}{
		{
			title: "account with no timezones",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`[]`)),
			},
			expected: []string{},
		},
		{
			title: "account with no timezones and empty server response body",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: nil,
		},
		{
			title: "account with some valid timezones",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`["tz"]`)),
			},
			expected: []string{"tz"},
		},
		{
			title:                "simulate remote call failure",
			response:             &http.Response{},
			forceClientSideError: true,
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
			httpClient := mock.NewHTTPClientMock(mock.ForceToFail(tC.forceClientSideError))
			client, err := New(
				WithBaseURL("https://base.com"),
				WithHTTPClient(httpClient),
				WithAPIKey("api_key"),
			)
			if err != nil {
				t.Errorf("Did not expect an error but received: '%v'", err)
				checkErrorType(t, err, true)
			}
			httpClient.SetResponse(fetchValidTimezonesPath, tC.response)
			actual, err := client.Accounts().Timezones()
			if err != nil {
				if !test.CheckError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceClientSideError)
			}
			if !reflect.DeepEqual(actual, tC.expected) {
				t.Errorf("Expected '%+v', actual: '%+v'", tC.expected, actual)
			}
		})
	}
}

func TestAccountsAPI_Now(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expected             time.Time
		expectedError        error
		parsingError         bool
	}{
		{
			title: "account without system time",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			},
			expected: time.Time{},
		},
		{
			title: "account without system time and empty server response body",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: time.Time{},
		},
		{
			title: "account with empty system time",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"SystemDate":""}`)),
			},
			expected: time.Time{},
		},
		{
			title: "account with parsable system time",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"SystemDate":"2020-06-12 16:19:00"}`)),
			},
			expected: time.Date(2020, 6, 12, 16, 19, 0, 0, time.UTC),
		},
		{
			title: "account with none parsable system time",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"SystemDate":"01/12/2006T16:19:00"}`)),
			},
			expectedError: newClientError(ErrCodeDataProcessing),
			expected:      time.Time{},
			parsingError:  true,
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
			httpClient := mock.NewHTTPClientMock(mock.ForceToFail(tC.forceHTTPClientError))
			client, err := New(
				WithBaseURL("https://base.com"),
				WithHTTPClient(httpClient),
				WithAPIKey("api_key"),
			)
			if err != nil {
				t.Errorf("Did not expect an error but received: '%v'", err)
				checkErrorType(t, err, true)
			}
			httpClient.SetResponse(fetchCurrentDatePath, tC.response)
			actual, err := client.Accounts().Now()
			if err != nil {
				if !test.CheckError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.parsingError && !tC.forceHTTPClientError)
			}
			if !reflect.DeepEqual(actual, tC.expected) {
				t.Errorf("Expected '%+v', actual: '%+v'", tC.expected, actual)
			}
		})
	}
}

func TestAccountsAPI_AddAdministrator(t *testing.T) {
	testCases := []struct {
		title                string
		forceClientSideError bool
		response             *http.Response
		expectedError        error
		input                accounts.Administrator
	}{
		{
			title: "receiving 200 from the server means success",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			},
		},
		{
			title:                "simulate remote call failure",
			response:             &http.Response{},
			forceClientSideError: true,
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
			httpClient := mock.NewHTTPClientMock(mock.ForceToFail(tC.forceClientSideError))
			client, err := New(
				WithBaseURL("https://base.com"),
				WithHTTPClient(httpClient),
				WithAPIKey("api_key"),
			)
			if err != nil {
				t.Errorf("Did not expect an error but received: '%v'", err)
				checkErrorType(t, err, true)
			}
			httpClient.SetResponse(administratorsPath, tC.response)
			err = client.Accounts().AddAdministrator(tC.input)
			if err != nil {
				if !test.CheckError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceClientSideError)
			}
		})
	}
}

func TestAccountsAPI_UpdateAdministrator(t *testing.T) {
	testCases := []struct {
		title                string
		forceClientSideError bool
		response             *http.Response
		expectedError        error
		input                accounts.Administrator
	}{
		{
			title: "receiving 200 from the server means success",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			},
		},
		{
			title:                "simulate remote call failure",
			response:             &http.Response{},
			forceClientSideError: true,
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
			httpClient := mock.NewHTTPClientMock(mock.ForceToFail(tC.forceClientSideError))
			client, err := New(
				WithBaseURL("https://base.com"),
				WithHTTPClient(httpClient),
				WithAPIKey("api_key"),
			)
			if err != nil {
				t.Errorf("Did not expect an error but received: '%v'", err)
				checkErrorType(t, err, true)
			}
			httpClient.SetResponse(administratorsPath, tC.response)
			err = client.Accounts().UpdateAdministrator("old+email@address.com", tC.input)
			if err != nil {
				if !test.CheckError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceClientSideError)
			}
		})
	}
}
