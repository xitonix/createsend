package createsend

import (
	"net/http"
	"strings"
	"time"
)

const (
	DefaultRetryCount = 3
	DefaultBaseURL    = "https://api.createsend.com/api/v3.2/"
)

type Option func(*Options)

type Options struct {
	retryCount int
	client     HTTPClient
	auth       *authentication
	baseURL    string
}

func defaultOptions() *Options {
	return &Options{
		retryCount: DefaultRetryCount,
		baseURL:    DefaultBaseURL,
		auth: &authentication{
			method: undefinedAuthentication,
		},
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func WithHTTPClient(client HTTPClient) Option {
	return func(options *Options) {
		options.client = client
	}
}

func WithBaseURL(url string) Option {
	return func(options *Options) {
		options.baseURL = url
	}
}

func WithAPIKey(apiKey string) Option {
	return func(options *Options) {
		options.auth = &authentication{
			token:  strings.TrimSpace(apiKey),
			method: apiKeyAuthentication,
		}
	}
}

func WithOAuthToken(token string) Option {
	return func(options *Options) {
		options.auth = &authentication{
			token:  strings.TrimSpace(token),
			method: oAuthAuthentication,
		}
	}
}

func WithRetryCount(count int) Option {
	return func(options *Options) {
		if count < 0 {
			count = 0
		}
		options.retryCount = count
	}
}
