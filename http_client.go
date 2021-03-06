package createsend

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	userAgentHeaderValue    = "createsend-go"
	userAgentHeaderKey      = "User-Agent"
	authenticationHeaderKey = "Authorization"
	httpsSchema             = "https://"
)

// HTTPClient is an interface for the internal HTTP client.
type HTTPClient interface {
	Do(request *http.Request) (*http.Response, error)
}

type httpClient struct {
	client  HTTPClient
	auth    *authentication
	baseURL *url.URL
	ctx     context.Context
}

func newHTTPClient(ctx context.Context, baseURL string, client HTTPClient, auth *authentication) (*httpClient, error) {
	if client == nil {
		return nil, newClientError(ErrCodeNilHTTPClient)
	}

	if err := auth.validate(); err != nil {
		return nil, err
	}

	baseURL = strings.ToLower(strings.TrimSpace(baseURL))
	if len(baseURL) == 0 {
		return nil, newClientError(ErrCodeEmptyURL)
	}

	// Enforce HTTPS
	baseURL = strings.Replace(baseURL, "http://", httpsSchema, 1)
	if !strings.HasPrefix(baseURL, httpsSchema) {
		baseURL = httpsSchema + baseURL
	}

	if baseURL == httpsSchema {
		return nil, newClientError(ErrCodeEmptyURL)
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, newWrappedClientError("Failed to parse the base URL", err, ErrCodeInvalidURL)
	}

	return &httpClient{
		client:  client,
		auth:    auth,
		baseURL: base,
		ctx:     ctx,
	}, nil
}

func (h *httpClient) Do(request *http.Request) (*http.Response, error) {
	request.Header.Add(userAgentHeaderKey, userAgentHeaderValue)
	h.auth.apply(request)
	return h.client.Do(request)
}

func (h *httpClient) Get(path string, result interface{}) error {
	return h.do(http.MethodGet, path, result, nil)
}

func (h *httpClient) Post(path string, result, body interface{}) error {
	return h.do(http.MethodPost, path, result, body)
}

func (h *httpClient) Put(path string, result, body interface{}) error {
	return h.do(http.MethodPut, path, result, body)
}

func (h *httpClient) Delete(path string) error {
	return h.do(http.MethodDelete, path, nil, nil)
}

func (h *httpClient) do(method, path string, result, body interface{}) error {
	request, err := h.newRequest(method, path, body)
	if err != nil {
		return err
	}
	response, err := h.Do(request)
	if err != nil {
		return newWrappedClientError("Failed to send the request", err, ErrCodeDataProcessing)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var csError Error
		err = json.NewDecoder(response.Body).Decode(&csError)
		if err != nil {
			return newWrappedClientError("Failed to decode the server error response", err, ErrCodeInvalidJSON)
		}
		csError.err = errors.New(csError.Message)
		return &csError
	}
	if result != nil {
		err = json.NewDecoder(response.Body).Decode(result)
		// io.EOF means that server returned an empty response
		// Successful HTTP status code with an empty response is Ok!
		if err != nil && !errors.Is(err, io.EOF) {
			return newWrappedClientError("Failed to decode the server response", err, ErrCodeInvalidJSON)
		}
	}
	return nil
}

func (h *httpClient) getFullURL(path string) (string, error) {
	path = strings.TrimSpace(path)
	if len(path) == 0 {
		return "", newClientError(ErrCodeEmptyURL)
	}
	rel, err := url.Parse(path)
	if err != nil {
		return "", newWrappedClientError("Failed to parse the request path", err, ErrCodeInvalidURL)
	}
	return h.baseURL.ResolveReference(rel).String(), nil
}

func (h *httpClient) newRequest(method, path string, body interface{}) (*http.Request, error) {
	fullURL, err := h.getFullURL(path)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if body != nil {
		err := json.NewEncoder(&buf).Encode(body)
		if err != nil {
			return nil, newWrappedClientError("Failed to serialise the request body", err, ErrCodeInvalidRequestBody)
		}
	}

	request, err := http.NewRequest(method, fullURL, &buf)
	if err != nil {
		return nil, newWrappedClientError("Failed to create the web request", err, ErrCodeDataProcessing)
	}

	return request.WithContext(h.ctx), nil
}
