package createsend

import (
	"net/http"
	"strings"
	"time"

	"github.com/xitonix/createsend/accounts"
)

const (
	DefaultBaseURL = "https://api.createsend.com/api/v3.2/"
)

type Option func(*Options)

type Options struct {
	client   HTTPClient
	auth     *authentication
	baseURL  string
	accounts accounts.API
}

func defaultOptions() *Options {
	return &Options{
		baseURL: DefaultBaseURL,
		auth: &authentication{
			method: undefinedAuthentication,
		},
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func WithAccountsAPI(api accounts.API) Option {
	return func(options *Options) {
		options.accounts = api
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
