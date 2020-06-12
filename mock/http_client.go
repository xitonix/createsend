package mock

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type result struct {
	count    int
	response *http.Response
}

type HTTPClientMock struct {
	options *Options
	lock    sync.Mutex
	calls   map[string]*result
}

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

func (h *HTTPClientMock) Do(request *http.Request) (*http.Response, error) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.options.callback != nil {
		h.options.callback(request)
	}
	call, ok := h.calls[request.URL.Path]
	if !ok {
		return nil, fmt.Errorf("no mocked response has been setup for %s. make sure you call SetResponse method first", request.URL.Path)
	}
	call.count++
	if h.options.forceToFail {
		return nil, ErrDeliberate
	}
	if h.options.fail != nil && h.options.fail(call.count, request) {
		return nil, ErrDeliberate
	}
	return call.response, nil
}

func (h *HTTPClientMock) Count(path string) int {
	call, ok := h.calls[path]
	if !ok {
		return -1
	}
	return call.count
}
