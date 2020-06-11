package mock

import "net/http"

const (
	DefaultRetryCount = 3
)

type CallFailurePredicate func(callCount int, request *http.Request) bool

type Option func(*Options)

type Options struct {
	forceToFail bool
	fail        CallFailurePredicate
	callback    func(*http.Request)
}

func defaultOptions() *Options {
	return &Options{}
}

func ForceToFail(force bool) Option {
	return func(options *Options) {
		options.forceToFail = force
	}
}

func WhenCalled(callback func(*http.Request)) Option {
	return func(options *Options) {
		options.callback = callback
	}
}

func FailIf(fail CallFailurePredicate) Option {
	return func(options *Options) {
		options.fail = fail
	}
}
