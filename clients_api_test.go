package createsend

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/xitonix/createsend/clients"
	"github.com/xitonix/createsend/mock"
)

func TestClientsAPI_Create(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expectedError        error
		input                clients.Client
		expectedResult       string
		oAuthAuthentication  bool
	}{
		{
			title: "successful execution",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`"client_id""`)),
			},
			expectedResult: "client_id",
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`"client_id"`)),
			},
			expectedResult:      "client_id",
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
			httpClient.SetResponse(clientsPath, tC.response)
			actual, err := client.Clients().Create(tC.input)
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceHTTPClientError)
			}
			if actual != tC.expectedResult {
				t.Errorf("Expected: %v, Actual: %v", tC.expectedResult, actual)
			}
		})
	}
}
