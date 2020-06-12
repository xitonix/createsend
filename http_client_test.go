package createsend

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
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

			if err != nil {
				checkErrorType(t, err, false)
			}

			if len(expectedHeaders) != len(request.Header) {
				t.Errorf("Expected Headers Count: %d, Actual: %d", len(expectedHeaders), len(request.Header))
			}

			for key, expectedValue := range expectedHeaders {
				if actual := request.Header.Get(key); actual != expectedValue {
					t.Errorf("Expected Header [%s]: '%s', Actual: '%s'", key, expectedValue, actual)
				}
			}
		})
	}
}

func TestNewHTTPSEnforcement(t *testing.T) {
	testCases := []struct {
		title    string
		baseURL  string
		expected string
	}{
		{
			title:    "no URL scheme is defined with host name only",
			baseURL:  "base.com",
			expected: "https://base.com",
		},
		{
			title:    "no URL scheme is defined with host name and path",
			baseURL:  "base.com/some/path",
			expected: "https://base.com/some/path",
		},
		{
			title:    "http should be replaced with https when only host name is defined",
			baseURL:  "http://base.com",
			expected: "https://base.com",
		},
		{
			title:    "http should be replaced with https when host name and path are defined",
			baseURL:  "http://base.com/some/path",
			expected: "https://base.com/some/path",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {
			auth := &authentication{
				method: apiKeyAuthentication,
				token:  "api_key",
			}
			client, err := newHTTPClient(tC.baseURL, mock.NewHTTPClientMock(), auth)
			if err != nil {
				checkErrorType(t, err, false)
			}
			actual := client.baseURL.String()
			if actual != tC.expected {
				t.Errorf("Expected base URL: %s, actual: %s", tC.expected, actual)
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
			if err != nil {
				checkErrorType(t, err, false)
			}
		})
	}
}

func TestNewRequest(t *testing.T) {
	const base = "https://base"
	testCases := []struct {
		title         string
		path          string
		method        string
		body          *bodyMock
		expectedError error
	}{
		{
			title:         "simulate invalid URL failure",
			method:        http.MethodGet,
			path:          "",
			body:          newBodyMock(false),
			expectedError: newClientError(ErrCodeEmptyURL),
		},
		{
			title:         "simulate body encoding failure",
			method:        http.MethodGet,
			path:          "/path",
			body:          newBodyMock(true),
			expectedError: newClientError(ErrCodeInvalidRequestBody),
		},
		{
			title:         "simulate request creation failure",
			method:        "こんにちは",
			path:          "/path",
			body:          newBodyMock(false),
			expectedError: newClientError(ErrCodeDataProcessingError),
		},
		{
			title:         "successful request creation",
			method:        http.MethodPut,
			path:          "/path",
			body:          newBodyMock(false),
			expectedError: nil,
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

			if err != nil {
				checkErrorType(t, err, false)
			}

			actual, err := client.newRequest(tC.method, tC.path, tC.body)
			if !test.CheckError(err, tC.expectedError) {
				t.Errorf("Expected '%v' error, but received: '%v'", tC.expectedError, err)
			}

			if err != nil {
				checkErrorType(t, err, false)
				return
			}

			expectedURL := base + tC.path

			if actual.URL.String() != expectedURL {
				t.Errorf("Expected request URL: %s, actual: %s", expectedURL, actual.URL.String())
			}

			if tC.method != actual.Method {
				t.Errorf("Expected %s HTTP method, actual: %s", tC.method, actual.Method)
			}

			checkRequestBody(t, actual.Body, tC.body)

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

			if err != nil {
				checkErrorType(t, err, false)
			}

			actual, err := client.getFullURL(tC.path)
			if !test.CheckError(err, tC.expectedError) {
				t.Errorf("Get Full URL: Expected '%v' error, but received: '%v'", tC.expectedError, err)
			}

			if err != nil {
				checkErrorType(t, err, false)
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
		expectServerError     bool
		forceRemoteCallToFail bool
	}{
		{
			title: "a successful server response with an empty body is acceptable",
			path:  "/path",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expectedError:  nil,
			expectedResult: &result{},
		},
		{
			title: "decoding valid json response from the server should not fail",
			path:  "/path",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":"d"}`)),
			},
			expectedError:  nil,
			expectedResult: &result{Data: "d"},
		},
		{
			title:                 "the client must report failure if the underlying http client failed",
			path:                  "/path",
			response:              &http.Response{},
			expectedError:         newClientError(ErrCodeDataProcessingError),
			expectedResult:        &result{},
			forceRemoteCallToFail: true,
		},
		{
			title: "the client must return a server error when the server returns a failure http status code",
			path:  "/path",
			response: &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":100}`)),
			},
			expectedError:     &Error{Code: 100},
			expectedResult:    &result{},
			expectServerError: true,
		},
		{
			title: "response status code 300 is considered a server side error",
			path:  "/path",
			response: &http.Response{
				StatusCode: 300,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":100}`)),
			},
			expectedError:     &Error{Code: 100},
			expectedResult:    &result{},
			expectServerError: true,
		},
		{
			title: "response status code greater than 300 is considered a server side error",
			path:  "/path",
			response: &http.Response{
				StatusCode: 301,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":100}`)),
			},
			expectedError:     &Error{Code: 100},
			expectedResult:    &result{},
			expectServerError: true,
		},
		{
			title: "response status code less than 200 is considered a server side error",
			path:  "/path",
			response: &http.Response{
				StatusCode: 199,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":100}`)),
			},
			expectedError:     &Error{Code: 100},
			expectedResult:    &result{},
			expectServerError: true,
		},
		{
			title: "response status code between 200 and 300 is considered a successful response",
			path:  "/path",
			response: &http.Response{
				StatusCode: 299,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":"d"}`)),
			},
			expectedError:  nil,
			expectedResult: &result{Data: "d"},
		},
		{
			title: "should fail to decode a server side error with invalid createsend error json content",
			path:  "/path",
			response: &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{I'm invalid Json'`)),
			},
			expectedError:  newClientError(ErrCodeInvalidJson),
			expectedResult: &result{},
		},
		{
			title: "should fail to decode the response with invalid json body",
			path:  "/path",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{I'm invalid Json'`)),
			},
			expectedError:  newClientError(ErrCodeInvalidJson),
			expectedResult: &result{},
		},
	}

	method := http.MethodGet
	expectedHeaderKeys := []string{
		userAgentHeaderKey,
		authenticationHeaderKey,
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {

			fullURL := base + tC.path

			onCall := func(request *http.Request) {
				if request.Method != method {
					t.Errorf("Expected HTTP Method: %s, Actual: %s", method, request.Method)
				}

				actualURL := request.URL.String()
				if actualURL != fullURL {
					t.Errorf("Expected URL: %s, Actual: %s", fullURL, actualURL)
				}

				if len(expectedHeaderKeys) != len(request.Header) {
					t.Errorf("Expected Headers Count: %d, Actual: %d", len(expectedHeaderKeys), len(request.Header))
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
			err = client.Get(fullURL, &actualResult)
			if !test.CheckError(err, tC.expectedError) {
				t.Errorf("Expected '%v' error, but received: '%v'", tC.expectedError, err)
			}

			if err != nil {
				checkErrorType(t, err, tC.expectServerError)
			}

			if !reflect.DeepEqual(&actualResult, tC.expectedResult) {
				t.Errorf("Expected result: %+v, actual: %+v", tC.expectedResult, actualResult)
			}

			count := httpClient.Count(tC.path)
			if count != 1 {
				t.Errorf("Expected number of calls: 1, actual: %d", count)
			}
		})
	}
}

func TestPost(t *testing.T) {
	const base = "https://base"
	type result struct {
		Data string `json:"data"`
	}

	testCases := []struct {
		title                       string
		path                        string
		response                    *http.Response
		expectedResult              *result
		body                        *bodyMock
		expectedError               error
		expectServerError           bool
		forceRemoteCallToFail       bool
		expectRemoteServerCallCount int
	}{
		{
			title: "a successful server response with an empty body is acceptable",
			path:  "/path",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expectedError:               nil,
			expectedResult:              &result{},
			expectRemoteServerCallCount: 1,
		},
		{
			title: "decoding valid json response from the server should not fail",
			path:  "/path",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":"d"}`)),
			},
			expectedError:               nil,
			expectedResult:              &result{Data: "d"},
			expectRemoteServerCallCount: 1,
		},
		{
			title:                       "the client must report failure if the underlying http client failed",
			path:                        "/path",
			response:                    &http.Response{},
			expectedError:               newClientError(ErrCodeDataProcessingError),
			expectedResult:              &result{},
			forceRemoteCallToFail:       true,
			expectRemoteServerCallCount: 1,
		},
		{
			title: "the client must return a server error when the server returns a failure http status code",
			path:  "/path",
			response: &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":100}`)),
			},
			expectedError:               &Error{Code: 100},
			expectedResult:              &result{},
			expectServerError:           true,
			expectRemoteServerCallCount: 1,
		},
		{
			title: "response status code 300 is considered a server side error",
			path:  "/path",
			response: &http.Response{
				StatusCode: 300,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":100}`)),
			},
			expectedError:               &Error{Code: 100},
			expectedResult:              &result{},
			expectServerError:           true,
			expectRemoteServerCallCount: 1,
		},
		{
			title: "response status code greater than 300 is considered a server side error",
			path:  "/path",
			response: &http.Response{
				StatusCode: 301,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":100}`)),
			},
			expectedError:               &Error{Code: 100},
			expectedResult:              &result{},
			expectServerError:           true,
			expectRemoteServerCallCount: 1,
		},
		{
			title: "response status code less than 200 is considered a server side error",
			path:  "/path",
			response: &http.Response{
				StatusCode: 199,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":100}`)),
			},
			expectedError:               &Error{Code: 100},
			expectedResult:              &result{},
			expectServerError:           true,
			expectRemoteServerCallCount: 1,
		},
		{
			title: "response status code between 200 and 300 is considered a successful response",
			path:  "/path",
			response: &http.Response{
				StatusCode: 299,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":"d"}`)),
			},
			expectedError:               nil,
			expectedResult:              &result{Data: "d"},
			expectRemoteServerCallCount: 1,
		},
		{
			title: "should fail to decode a server side error with invalid createsend error json content",
			path:  "/path",
			response: &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{I'm invalid Json'`)),
			},
			expectedError:               newClientError(ErrCodeInvalidJson),
			expectedResult:              &result{},
			expectRemoteServerCallCount: 1,
		},
		{
			title: "none empty valid body should be marshalled into the request",
			path:  "/path",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":"d"}`)),
			},
			expectedError:               nil,
			body:                        newBodyMock(false),
			expectedResult:              &result{Data: "d"},
			expectRemoteServerCallCount: 1,
		},
		{
			title: "request should not be sent to the server if encoding of request body failed",
			path:  "/path",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":"d"}`)),
			},
			expectedError:               newClientError(ErrCodeInvalidRequestBody),
			body:                        newBodyMock(true),
			expectedResult:              &result{},
			expectRemoteServerCallCount: 0,
		},
		{
			title: "should fail to decode the response with invalid json body",
			path:  "/path",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{I'm invalid Json'`)),
			},
			expectedError:               newClientError(ErrCodeInvalidJson),
			expectedResult:              &result{},
			expectRemoteServerCallCount: 1,
		},
	}

	method := http.MethodPost
	expectedHeaderKeys := []string{
		userAgentHeaderKey,
		authenticationHeaderKey,
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {

			fullURL := base + tC.path

			onCall := func(request *http.Request) {
				if request.Method != method {
					t.Errorf("Expected HTTP Method: %s, Actual: %s", method, request.Method)
				}

				actualURL := request.URL.String()
				if actualURL != fullURL {
					t.Errorf("Expected URL: %s, Actual: %s", fullURL, actualURL)
				}

				if len(expectedHeaderKeys) != len(request.Header) {
					t.Errorf("Expected Headers Count: %d, Actual: %d", len(expectedHeaderKeys), len(request.Header))
				}

				for _, key := range expectedHeaderKeys {
					if _, ok := request.Header[key]; !ok {
						t.Errorf("Expected Header [%s] was not found", key)
					}
				}

				if tC.body != nil {
					checkRequestBody(t, request.Body, tC.body)
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
			err = client.Post(fullURL, &actualResult, tC.body)
			if !test.CheckError(err, tC.expectedError) {
				t.Errorf("Expected '%v' error, but received: '%v'", tC.expectedError, err)
			}

			if err != nil {
				checkErrorType(t, err, tC.expectServerError)
			}

			if !reflect.DeepEqual(&actualResult, tC.expectedResult) {
				t.Errorf("Expected result: %+v, actual: %+v", tC.expectedResult, actualResult)
			}

			count := httpClient.Count(tC.path)
			if tC.expectRemoteServerCallCount != count {
				t.Errorf("Expected number of calls to remote servers: %d, actual: %d", tC.expectRemoteServerCallCount, count)
			}
		})
	}
}

func TestPut(t *testing.T) {
	const base = "https://base"
	type result struct {
		Data string `json:"data"`
	}

	testCases := []struct {
		title                       string
		path                        string
		response                    *http.Response
		expectedResult              *result
		body                        *bodyMock
		expectedError               error
		expectServerError           bool
		forceRemoteCallToFail       bool
		expectRemoteServerCallCount int
	}{
		{
			title: "a successful server response with an empty body is acceptable",
			path:  "/path",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expectedError:               nil,
			expectedResult:              &result{},
			expectRemoteServerCallCount: 1,
		},
		{
			title: "decoding valid json response from the server should not fail",
			path:  "/path",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":"d"}`)),
			},
			expectedError:               nil,
			expectedResult:              &result{Data: "d"},
			expectRemoteServerCallCount: 1,
		},
		{
			title:                       "the client must report failure if the underlying http client failed",
			path:                        "/path",
			response:                    &http.Response{},
			expectedError:               newClientError(ErrCodeDataProcessingError),
			expectedResult:              &result{},
			forceRemoteCallToFail:       true,
			expectRemoteServerCallCount: 1,
		},
		{
			title: "the client must return a server error when the server returns a failure http status code",
			path:  "/path",
			response: &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":100}`)),
			},
			expectedError:               &Error{Code: 100},
			expectedResult:              &result{},
			expectServerError:           true,
			expectRemoteServerCallCount: 1,
		},
		{
			title: "response status code 300 is considered a server side error",
			path:  "/path",
			response: &http.Response{
				StatusCode: 300,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":100}`)),
			},
			expectedError:               &Error{Code: 100},
			expectedResult:              &result{},
			expectServerError:           true,
			expectRemoteServerCallCount: 1,
		},
		{
			title: "response status code greater than 300 is considered a server side error",
			path:  "/path",
			response: &http.Response{
				StatusCode: 301,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":100}`)),
			},
			expectedError:               &Error{Code: 100},
			expectedResult:              &result{},
			expectServerError:           true,
			expectRemoteServerCallCount: 1,
		},
		{
			title: "response status code less than 200 is considered a server side error",
			path:  "/path",
			response: &http.Response{
				StatusCode: 199,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":100}`)),
			},
			expectedError:               &Error{Code: 100},
			expectedResult:              &result{},
			expectServerError:           true,
			expectRemoteServerCallCount: 1,
		},
		{
			title: "response status code between 200 and 300 is considered a successful response",
			path:  "/path",
			response: &http.Response{
				StatusCode: 299,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":"d"}`)),
			},
			expectedError:               nil,
			expectedResult:              &result{Data: "d"},
			expectRemoteServerCallCount: 1,
		},
		{
			title: "should fail to decode a server side error with invalid createsend error json content",
			path:  "/path",
			response: &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{I'm invalid Json'`)),
			},
			expectedError:               newClientError(ErrCodeInvalidJson),
			expectedResult:              &result{},
			expectRemoteServerCallCount: 1,
		},
		{
			title: "none empty valid body should be marshalled into the request",
			path:  "/path",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":"d"}`)),
			},
			expectedError:               nil,
			body:                        newBodyMock(false),
			expectedResult:              &result{Data: "d"},
			expectRemoteServerCallCount: 1,
		},
		{
			title: "request should not be sent to the server if encoding of request body failed",
			path:  "/path",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":"d"}`)),
			},
			expectedError:               newClientError(ErrCodeInvalidRequestBody),
			body:                        newBodyMock(true),
			expectedResult:              &result{},
			expectRemoteServerCallCount: 0,
		},
		{
			title: "should fail to decode the response with invalid json body",
			path:  "/path",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{I'm invalid Json'`)),
			},
			expectedError:               newClientError(ErrCodeInvalidJson),
			expectedResult:              &result{},
			expectRemoteServerCallCount: 1,
		},
	}

	method := http.MethodPut
	expectedHeaderKeys := []string{
		userAgentHeaderKey,
		authenticationHeaderKey,
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {

			fullURL := base + tC.path

			onCall := func(request *http.Request) {
				if request.Method != method {
					t.Errorf("Expected HTTP Method: %s, Actual: %s", method, request.Method)
				}

				actualURL := request.URL.String()
				if actualURL != fullURL {
					t.Errorf("Expected URL: %s, Actual: %s", fullURL, actualURL)
				}

				if len(expectedHeaderKeys) != len(request.Header) {
					t.Errorf("Expected Headers Count: %d, Actual: %d", len(expectedHeaderKeys), len(request.Header))
				}

				for _, key := range expectedHeaderKeys {
					if _, ok := request.Header[key]; !ok {
						t.Errorf("Expected Header [%s] was not found", key)
					}
				}

				if tC.body != nil {
					checkRequestBody(t, request.Body, tC.body)
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
			err = client.Put(fullURL, &actualResult, tC.body)
			if !test.CheckError(err, tC.expectedError) {
				t.Errorf("Expected '%v' error, but received: '%v'", tC.expectedError, err)
			}

			if err != nil {
				checkErrorType(t, err, tC.expectServerError)
			}

			if !reflect.DeepEqual(&actualResult, tC.expectedResult) {
				t.Errorf("Expected result: %+v, actual: %+v", tC.expectedResult, actualResult)
			}

			count := httpClient.Count(tC.path)
			if tC.expectRemoteServerCallCount != count {
				t.Errorf("Expected number of calls to remote servers: %d, actual: %d", tC.expectRemoteServerCallCount, count)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	const base = "https://base"
	testCases := []struct {
		title                 string
		path                  string
		response              *http.Response
		expectedError         error
		expectServerError     bool
		forceRemoteCallToFail bool
	}{
		{
			title:                 "the client must report failure if the underlying http client failed",
			path:                  "/path",
			response:              &http.Response{},
			expectedError:         newClientError(ErrCodeDataProcessingError),
			forceRemoteCallToFail: true,
		},
		{
			title: "the client must return a server error when the server returns a failure http status code",
			path:  "/path",
			response: &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":100}`)),
			},
			expectedError:     &Error{Code: 100},
			expectServerError: true,
		},
		{
			title: "response status code 300 is considered a server side error",
			path:  "/path",
			response: &http.Response{
				StatusCode: 300,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":100}`)),
			},
			expectedError:     &Error{Code: 100},
			expectServerError: true,
		},
		{
			title: "response status code greater than 300 is considered a server side error",
			path:  "/path",
			response: &http.Response{
				StatusCode: 301,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":100}`)),
			},
			expectedError:     &Error{Code: 100},
			expectServerError: true,
		},
		{
			title: "response status code less than 200 is considered a server side error",
			path:  "/path",
			response: &http.Response{
				StatusCode: 199,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":100}`)),
			},
			expectedError:     &Error{Code: 100},
			expectServerError: true,
		},
		{
			title: "response status code between 200 and 300 is considered a successful response",
			path:  "/path",
			response: &http.Response{
				StatusCode: 299,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":"d"}`)),
			},
			expectedError: nil,
		},
		{
			title: "should fail to decode a server side error with invalid createsend error json content",
			path:  "/path",
			response: &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{I'm invalid Json'`)),
			},
			expectedError: newClientError(ErrCodeInvalidJson),
		},
	}

	method := http.MethodDelete
	expectedHeaderKeys := []string{
		userAgentHeaderKey,
		authenticationHeaderKey,
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {

			fullURL := base + tC.path

			onCall := func(request *http.Request) {
				if request.Method != method {
					t.Errorf("Expected HTTP Method: %s, Actual: %s", method, request.Method)
				}

				actualURL := request.URL.String()
				if actualURL != fullURL {
					t.Errorf("Expected URL: %s, Actual: %s", fullURL, actualURL)
				}

				if len(expectedHeaderKeys) != len(request.Header) {
					t.Errorf("Expected Headers Count: %d, Actual: %d", len(expectedHeaderKeys), len(request.Header))
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

			err = client.Delete(fullURL)
			if !test.CheckError(err, tC.expectedError) {
				t.Errorf("Expected '%v' error, but received: '%v'", tC.expectedError, err)
			}

			if err != nil {
				checkErrorType(t, err, tC.expectServerError)
			}

			count := httpClient.Count(tC.path)
			if count != 1 {
				t.Errorf("Expected number of calls: 1, actual: %d", count)
			}
		})
	}
}

func checkRequestBody(t *testing.T, actual io.ReadCloser, expected *bodyMock) {
	t.Helper()
	defer actual.Close()
	b, err := ioutil.ReadAll(actual)
	if err != nil {
		t.Errorf("Did not expect an error but received: '%v'", err)
	}
	var body bodyMock
	err = json.Unmarshal(b, &body)
	if err != nil {
		t.Errorf("Did not expect an error but received: '%v'", err)
	}

	if expected == nil {
		if len(b) > 0 {
			t.Errorf("The request body should have been empty")
		}
		return
	}
	if expected.Value != body.Value {
		t.Errorf("Expected Request Body Value: %s, Actual: %s", expected.Value, body.Value)
	}
}

func checkErrorType(t *testing.T, err error, expectServerError bool) {
	t.Helper()
	var csErr *Error
	ok := errors.As(err, &csErr)
	if !ok {
		t.Error("We should always return a custom createsend Error type")
	}
	if csErr.IsFromServer() != expectServerError {
		t.Errorf("Expected server error: %v, actual: %v", expectServerError, csErr.IsFromServer())
	}
}
