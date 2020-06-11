package createsend

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/xitonix/createsend/internal/test"
	"github.com/xitonix/createsend/mock"
)

func TestHeaders(t *testing.T) {
	testCases := []struct {
		title                        string
		auth                         *authentication
		expectedAuthenticationHeader string
		request                      *http.Request
	}{
		{
			title: "with api key authentication",
			auth: &authentication{
				token:  "api_key",
				method: apiKeyAuthentication,
			},
			expectedAuthenticationHeader: "Basic YXBpX2tleTphcGlfa2V5",
		},
		{
			title: "with oauth token authentication",
			auth: &authentication{
				token:  "token",
				method: oAuthAuthentication,
			},
			expectedAuthenticationHeader: "Bearer token",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {
			httpClient := mock.NewHTTPClientMock()
			client, err := newHTTPClient("https://base", httpClient, tC.auth)

			if err != nil {
				t.Errorf("Client Creation: Did not expect to receive an error, but received: '%s'", err)
			}

			expectedHeaders := map[string]string{
				userAgentHeaderKey:      userAgentHeaderValue,
				authenticationHeaderKey: tC.expectedAuthenticationHeader,
			}
			url := "/dummy"
			httpClient.SetResponse(url, &http.Response{})
			request, _ := http.NewRequest(http.MethodGet, url, nil)
			_, err = client.Do(request)
			if err != nil {
				t.Errorf("Do: did not expect an error, but received '%s'", err)
			}

			if len(expectedHeaders) != len(request.Header) {
				t.Errorf("Expected Headers Length: %d, Actual: %d", len(expectedHeaders), len(request.Header))
			}

			for key, expectedValue := range expectedHeaders {
				if actual := request.Header.Get(key); actual != expectedValue {
					t.Errorf("Expected Header [%s]: '%s', Actual: '%s'", key, expectedValue, actual)
				}
			}
		})
	}
}

func TestNewHTTPClient(t *testing.T) {
	testCases := []struct {
		title         string
		auth          *authentication
		expectedError error
		client        HTTPClient
		baseURL       string
	}{
		{
			title:         "nil authentication",
			auth:          nil,
			expectedError: newClientError(ErrCodeAuthenticationNotSet),
			client:        mock.NewHTTPClientMock(),
			baseURL:       "https://base",
		},
		{
			title: "empty base URL",
			auth: &authentication{
				method: apiKeyAuthentication,
				token:  "api_key",
			},
			client:        mock.NewHTTPClientMock(),
			expectedError: newClientError(ErrCodeEmptyURL),
			baseURL:       "",
		},
		{
			title: "whitespace base URL",
			auth: &authentication{
				method: apiKeyAuthentication,
				token:  "api_key",
			},
			client:        mock.NewHTTPClientMock(),
			expectedError: newClientError(ErrCodeEmptyURL),
			baseURL:       "   ",
		},
		{
			title: "http schema only base URL",
			auth: &authentication{
				method: apiKeyAuthentication,
				token:  "api_key",
			},
			client:        mock.NewHTTPClientMock(),
			expectedError: newClientError(ErrCodeEmptyURL),
			baseURL:       "http://",
		},
		{
			title: "https schema only base URL",
			auth: &authentication{
				method: apiKeyAuthentication,
				token:  "api_key",
			},
			client:        mock.NewHTTPClientMock(),
			expectedError: newClientError(ErrCodeEmptyURL),
			baseURL:       "https://",
		},
		{
			title: "invalid base URL",
			auth: &authentication{
				method: apiKeyAuthentication,
				token:  "api_key",
			},
			client:        mock.NewHTTPClientMock(),
			expectedError: newClientError(ErrCodeInvalidURL),
			baseURL:       "%",
		},
		{
			title: "undefined auth method",
			auth: &authentication{
				token: "token",
			},
			expectedError: newClientError(ErrCodeAuthenticationNotSet),
			client:        mock.NewHTTPClientMock(),
			baseURL:       "https://base",
		},
		{
			title: "empty api key method",
			auth: &authentication{
				method: apiKeyAuthentication,
				token:  "",
			},
			expectedError: newClientError(ErrCodeEmptyAPIKey),
			client:        mock.NewHTTPClientMock(),
			baseURL:       "https://base",
		},
		{
			title: "nil http client",
			auth: &authentication{
				method: apiKeyAuthentication,
				token:  "",
			},
			expectedError: newClientError(ErrCodeNilHTTPClient),
			client:        nil,
			baseURL:       "https://base",
		},
		{
			title: "empty oauth token",
			auth: &authentication{
				method: oAuthAuthentication,
				token:  "",
			},
			expectedError: newClientError(ErrCodeEmptyOAuthToken),
			client:        mock.NewHTTPClientMock(),
			baseURL:       "https://base",
		},
		{
			title: "with api key authentication",
			auth: &authentication{
				token:  "api_key",
				method: apiKeyAuthentication,
			},
			client:  mock.NewHTTPClientMock(),
			baseURL: "https://base",
		},
		{
			title: "with oauth token authentication",
			auth: &authentication{
				token:  "token",
				method: oAuthAuthentication,
			},
			client:  mock.NewHTTPClientMock(),
			baseURL: "https://base",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {
			_, err := newHTTPClient(tC.baseURL, tC.client, tC.auth)
			if !test.CheckError(err, tC.expectedError) {
				t.Errorf("Client Creation: Expected '%v' error, but received: '%v'", tC.expectedError, err)
			}
		})
	}
}

func TestGetFullURL(t *testing.T) {
	const base = "https://base"
	testCases := []struct {
		title         string
		path          string
		expectedError error
		expectedURL   string
	}{
		{
			title:         "with empty path",
			expectedError: newClientError(ErrCodeEmptyURL),
		},
		{
			title:         "whitespace path",
			expectedError: newClientError(ErrCodeEmptyURL),
		},
		{
			title:         "invalid path",
			path:          "%",
			expectedError: newClientError(ErrCodeInvalidURL),
		},
		{
			title:       "valid path",
			path:        "/path",
			expectedURL: base + "/path",
		},
		{
			title:       "valid path with extension",
			path:        "/path/resource.json",
			expectedURL: base + "/path/resource.json",
		},
		{
			title:       "valid path with query string parameters",
			path:        "/path?a=1&b=true",
			expectedURL: base + "/path?a=1&b=true",
		},
		{
			title:       "valid path with extension and query string parameters",
			path:        "/path/resource.json?a=1&b=true",
			expectedURL: base + "/path/resource.json?a=1&b=true",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {
			httpClient := mock.NewHTTPClientMock()
			auth := &authentication{
				token:  "api_key",
				method: apiKeyAuthentication,
			}
			client, err := newHTTPClient(base, httpClient, auth)
			if err != nil {
				t.Errorf("Client Creation: Did not expect to receive an error, but received: '%s'", err)
			}

			actual, err := client.getFullURL(tC.path)
			if !test.CheckError(err, tC.expectedError) {
				t.Errorf("Get Full URL: Expected '%v' error, but received: '%v'", tC.expectedError, err)
			}

			if actual != tC.expectedURL {
				t.Errorf("Get Full URL: Expected URL: '%s', actual: '%s'", tC.expectedURL, actual)
			}
		})
	}
}

func TestGet(t *testing.T) {
	const base = "https://base"
	type result struct {
		Data string `json:"data"`
	}
	testCases := []struct {
		title                 string
		path                  string
		response              *http.Response
		expectedResult        *result
		expectedError         error
		forceRemoteCallToFail bool
	}{
		{
			title: "empty body",
			path:  "/path",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
			},
			expectedError:  newClientError(ErrCodeInvalidJson),
			expectedResult: &result{},
		},
		{
			title: "none empty body",
			path:  "/path",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":"d"}`)),
			},
			expectedError:  nil,
			expectedResult: &result{Data: "d"},
		},
		{
			title:                 "fail to call the server",
			path:                  "/path",
			response:              &http.Response{},
			expectedError:         newClientError(ErrCodeUnknown),
			expectedResult:        &result{},
			forceRemoteCallToFail: true,
		},
		{
			title: "server side error with valid createsend error body",
			path:  "/path",
			response: &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":100}`)),
			},
			expectedError:  &Error{Code: 100},
			expectedResult: &result{},
		},
		{
			title: "server side error with invalid createsend error body",
			path:  "/path",
			response: &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
			},
			expectedError:  newClientError(ErrCodeInvalidJson),
			expectedResult: &result{},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {
			expectedHeaderKeys := []string{
				userAgentHeaderKey,
				authenticationHeaderKey,
			}

			onCall := func(request *http.Request) {
				if len(expectedHeaderKeys) != len(request.Header) {
					t.Errorf("Expected Headers Length: %d, Actual: %d", len(expectedHeaderKeys), len(request.Header))
				}

				for _, key := range expectedHeaderKeys {
					if _, ok := request.Header[key]; !ok {
						t.Errorf("Expected Header [%s] was not found", key)
					}
				}
			}
			httpClient := mock.NewHTTPClientMock(mock.WhenCalled(onCall), mock.ForceToFail(tC.forceRemoteCallToFail))
			auth := &authentication{
				token:  "api_key",
				method: apiKeyAuthentication,
			}

			client, err := newHTTPClient(base, httpClient, auth)
			if err != nil {
				t.Errorf("Client Creation: Did not expect to receive an error, but received: '%s'", err)
			}

			httpClient.SetResponse(tC.path, tC.response)

			var actualResult result
			err = client.Get(base+tC.path, &actualResult)
			if !test.CheckError(err, tC.expectedError) {
				t.Errorf("Get: Expected '%v' error, but received: '%v'", tC.expectedError, err)
			}

			if !reflect.DeepEqual(&actualResult, tC.expectedResult) {
				t.Errorf("Get: Expected result: %+v, actual: %+v", tC.expectedResult, actualResult)
			}

			count := httpClient.Count(tC.path)
			if count != 1 {
				t.Errorf("Get: Expected number of calls: 1, actual: %d", count)
			}
		})
	}
}
