package mock

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type result struct {
	count    int
	response *http.Response
}

// HTTPClientMock represents a mocked HTTP client.
type HTTPClientMock struct {
	options   *Options
	lock      sync.Mutex
	calls     map[string]*result
	requested *url.URL
}

// NewHTTPClientMock creates a new instance of a mocked HTTP client.
func NewHTTPClientMock(options ...Option) *HTTPClientMock {
	opts := defaultOptions()
	for _, op := range options {
		op(opts)
	}
	return &HTTPClientMock{
		options: opts,
		calls:   make(map[string]*result),
	}
}

// SetResponse sets the response you expect to be returned from the server once the specified path is hit.
func (h *HTTPClientMock) SetResponse(path string, response *http.Response) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	h.calls[path] = &result{
		count:    0,
		response: response,
	}
}

// Do sends an HTTP request to the mocked server.
func (h *HTTPClientMock) Do(request *http.Request) (*http.Response, error) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.options.callback != nil {
		h.options.callback(request)
	}
	h.requested = request.URL
	call, ok := h.calls[request.URL.Path]
	if !ok {
		return nil, fmt.Errorf("no mocked response has been setup for %s. make sure you call SetResponse method first", request.URL.Path)
	}
	call.count++
	if h.options.forceToFail {
		return nil, ErrDeliberate
	}
	return call.response, nil
}

// LastRequest returns the URL of the last requested resource.
func (h *HTTPClientMock) LastRequest() *url.URL {
	h.lock.Lock()
	defer h.lock.Unlock()
	return h.requested
}

// Count returns number of times the specified path was hit.
func (h *HTTPClientMock) Count(path string) int {
	call, ok := h.calls[path]
	if !ok {
		return -1
	}
	return call.count
}
