package mock

import "net/http"

// Option represents an optional configuration function for mocked clients.
type Option func(*Options)

// Options mocked client configurations.
type Options struct {
	forceToFail bool
	callback    func(*http.Request)
}

func defaultOptions() *Options {
	return &Options{}
}

// ForceToFail if enabled, it forces the mocked HTTP requests to fail.
func ForceToFail(force bool) Option {
	return func(options *Options) {
		options.forceToFail = force
	}
}

// WhenCalled sets a callback which will be called once an HTTP request is being sent to the server.
func WhenCalled(callback func(*http.Request)) Option {
	return func(options *Options) {
		options.callback = callback
	}
}
