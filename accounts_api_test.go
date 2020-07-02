package createsend

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/xitonix/createsend/accounts"
	"github.com/xitonix/createsend/mock"
)

func TestAccountsAPI_Clients(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expected             []*accounts.Client
		expectedError        error
		oAuthAuthentication  bool
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
				Id:   "id",
				Name: "name",
			}},
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`[{"ClientID":"id", "Name":"name"}]`)),
			},
			expected: []*accounts.Client{{
				Id:   "id",
				Name: "name",
			}},
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
			httpClient.SetResponse(listClientsPath, tC.response)
			actual, err := client.Accounts().Clients()
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

func TestAccountsAPI_Billing(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expected             *accounts.Billing
		expectedError        error
		oAuthAuthentication  bool
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
			expected: &accounts.Billing{},
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
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Credits":100}`)),
			},
			expected: &accounts.Billing{
				Credits: 100,
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
			httpClient.SetResponse(fetchBillingDetailsPath, tC.response)
			actual, err := client.Accounts().Billing()
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

func TestAccountsAPI_Countries(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expected             []string
		expectedError        error
		oAuthAuthentication  bool
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
			expected: []string{},
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
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`["country"]`)),
			},
			expected:            []string{"country"},
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
			httpClient.SetResponse(fetchValidCountriesPath, tC.response)
			actual, err := client.Accounts().Countries()
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

func TestAccountsAPI_Timezones(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expected             []string
		expectedError        error
		oAuthAuthentication  bool
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
			expected: []string{},
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
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`["tz"]`)),
			},
			expected:            []string{"tz"},
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
			httpClient.SetResponse(fetchValidTimezonesPath, tC.response)
			actual, err := client.Accounts().Timezones()
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

func TestAccountsAPI_Now(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expected             time.Time
		expectedError        error
		parsingError         bool
		oAuthAuthentication  bool
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
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"SystemDate":"2020-06-12 16:19:00"}`)),
			},
			expected:            time.Date(2020, 6, 12, 16, 19, 0, 0, time.UTC),
			oAuthAuthentication: true,
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
			client, httpClient := createClient(t, tC.oAuthAuthentication, tC.forceHTTPClientError)
			httpClient.SetResponse(fetchCurrentDatePath, tC.response)
			actual, err := client.Accounts().Now()
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.parsingError && !tC.forceHTTPClientError)
			}
			if diff := cmp.Diff(tC.expected, actual); diff != "" {
				t.Errorf("Expectations failed (-expected +actual):\n%s", diff)
			}
		})
	}
}

func TestAccountsAPI_AddAdministrator(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expectedError        error
		input                accounts.Administrator
		oAuthAuthentication  bool
	}{
		{
			title: "receiving 200 from the server means success",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			},
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
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
			httpClient.SetResponse(administratorsPath, tC.response)
			err := client.Accounts().AddAdministrator(tC.input)
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceHTTPClientError)
			}
		})
	}
}

func TestAccountsAPI_UpdateAdministrator(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expectedError        error
		oAuthAuthentication  bool
	}{
		{
			title: "receiving 200 from the server means success",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			},
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
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
			httpClient.SetResponse(administratorsPath, tC.response)
			err := client.Accounts().UpdateAdministrator("old+email@address.com", accounts.Administrator{})
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceHTTPClientError)
			}
		})
	}
}

func TestAccountsAPI_Administrators(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expected             []*accounts.AdministratorDetails
		expectedError        error
		oAuthAuthentication  bool
	}{
		{
			title: "account with no admins",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`[]`)),
			},
			expected: []*accounts.AdministratorDetails{},
		},
		{
			title: "account with no admins and empty server response body",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: []*accounts.AdministratorDetails{},
		},
		{
			title: "account with admins",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`[{"EmailAddress":"e@d.com", "Name":"name", "Status":"status"}]`)),
			},
			expected: []*accounts.AdministratorDetails{{
				Administrator: accounts.Administrator{
					EmailAddress: "e@d.com",
					Name:         "name",
				},
				Status: "status",
			}},
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`[{"EmailAddress":"e@d.com", "Name":"name", "Status":"status"}]`)),
			},
			expected: []*accounts.AdministratorDetails{{
				Administrator: accounts.Administrator{
					EmailAddress: "e@d.com",
					Name:         "name",
				},
				Status: "status",
			}},
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
			httpClient.SetResponse(administratorsPath, tC.response)
			actual, err := client.Accounts().Administrators()
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

func TestAccountsAPI_Administrator(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expected             *accounts.AdministratorDetails
		expectedError        error
		oAuthAuthentication  bool
	}{
		{
			title: "account with no admin",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			},
			expected: &accounts.AdministratorDetails{},
		},
		{
			title: "account with no admins and empty server response body",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: &accounts.AdministratorDetails{},
		},
		{
			title: "administrator found",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"EmailAddress":"e@d.com", "Name":"name", "Status":"status"}`)),
			},
			expected: &accounts.AdministratorDetails{
				Administrator: accounts.Administrator{
					EmailAddress: "e@d.com",
					Name:         "name",
				},
				Status: "status",
			},
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"EmailAddress":"e@d.com", "Name":"name", "Status":"status"}`)),
			},
			expected: &accounts.AdministratorDetails{
				Administrator: accounts.Administrator{
					EmailAddress: "e@d.com",
					Name:         "name",
				},
				Status: "status",
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
			httpClient.SetResponse(administratorsPath, tC.response)
			actual, err := client.Accounts().Administrator("e@d.com")
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

func TestAccountsAPI_DeleteAdministrator(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expectedError        error
		oAuthAuthentication  bool
	}{
		{
			title: "successful deletion",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
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
			httpClient.SetResponse(administratorsPath, tC.response)
			err := client.Accounts().DeleteAdministrator("e@d.com")
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceHTTPClientError)
			}
		})
	}
}

func TestAccountsAPI_PrimaryContact(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expected             string
		expectedError        error
		oAuthAuthentication  bool
	}{
		{
			title: "account with no primary contact",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			},
			expected: "",
		},
		{
			title: "account with no primary contact and empty server response body",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: "",
		},
		{
			title: "primary contact found",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"EmailAddress":"e@d.com"}`)),
			},
			expected: "e@d.com",
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"EmailAddress":"e@d.com"}`)),
			},
			expected:            "e@d.com",
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
			httpClient.SetResponse(primaryContactPath, tC.response)
			actual, err := client.Accounts().PrimaryContact()
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

func TestAccountsAPI_SetAsPrimaryContact(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expectedError        error
		oAuthAuthentication  bool
	}{
		{
			title: "successful execution",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"EmailAddress":"e@d.com"}`)),
			},
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"EmailAddress":"e@d.com"}`)),
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
			httpClient.SetResponse(primaryContactPath, tC.response)
			err := client.Accounts().SetAsPrimaryContact("e@d.com")
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceHTTPClientError)
			}
		})
	}
}

func TestAccountsAPI_NewEmbeddedSession(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		input                accounts.EmbeddedSession
		expected             string
		expectedError        error
		oAuthAuthentication  bool
	}{
		{
			title: "empty json object response",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			},
			expected: "",
		},
		{
			title: "empty server response body",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: "",
		},
		{
			title: "successful execution",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"SessionUrl":"session_url"}`)),
			},
			input: accounts.EmbeddedSession{
				EmailAddress: "e@d.com",
				Chrome:       "chrome",
				URL:          "url",
				IntegratorId: "integration_id",
				ClientId:     "client_id",
			},
			expected: "session_url",
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"SessionUrl":"session_url"}`)),
			},
			input: accounts.EmbeddedSession{
				EmailAddress: "e@d.com",
				Chrome:       "chrome",
				URL:          "url",
				IntegratorId: "integration_id",
				ClientId:     "client_id",
			},
			expected:            "session_url",
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
			httpClient.SetResponse(externalSessionPath, tC.response)
			actual, err := client.Accounts().NewEmbeddedSession(tC.input)
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceHTTPClientError)
			}
			if actual != tC.expected {
				t.Errorf("Expected '%+v', actual: '%+v'", tC.expected, actual)
			}
		})
	}
}
