package createsend

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/xitonix/createsend/accounts"
	"github.com/xitonix/createsend/campaigns"
	"github.com/xitonix/createsend/clients"
	"github.com/xitonix/createsend/transactional"
)

const (
	// DefaultBaseURL the default API base URL.
	DefaultBaseURL = "https://api.createsend.com/api/v3.2/"
)

// Option represents an optional client configuration function.
type Option func(*Options)

// Options client configurations.
type Options struct {
	client        HTTPClient
	auth          *authentication
	baseURL       string
	accounts      accounts.API
	clients       clients.API
	transactional transactional.API
	campaigns     campaigns.API
	ctx           context.Context
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
		ctx: context.Background(),
	}
}

// WithClientsAPI overrides the internal object for accessing Clients API.
//
// You can override the API to mock out Clients API methods altogether.
func WithClientsAPI(api clients.API) Option {
	return func(options *Options) {
		options.clients = api
	}
}

// WithAccountsAPI overrides the internal object for accessing Accounts API.
//
// You can override the API to mock out Accounts API methods altogether.
func WithAccountsAPI(api accounts.API) Option {
	return func(options *Options) {
		options.accounts = api
	}
}

// WithTransactionalAPI overrides the internal object for accessing Transactional API.
//
// You can override the API to mock out Transactional API methods altogether.
func WithTransactionalAPI(api transactional.API) Option {
	return func(options *Options) {
		options.transactional = api
	}
}

// WithCampaignsAPI overrides the internal object for accessing Campaign API.
//
// You can override the API to mock out Campaign API methods altogether.
func WithCampaignsAPI(api campaigns.API) Option {
	return func(options *Options) {
		options.campaigns = api
	}
}

// WithContext sets the context for all the HTTP requests.
func WithContext(ctx context.Context) Option {
	return func(options *Options) {
		if ctx == nil {
			ctx = context.Background()
		}
		options.ctx = ctx
	}
}

// WithHTTPClient sets the internal HTTP client.
func WithHTTPClient(client HTTPClient) Option {
	return func(options *Options) {
		options.client = client
	}
}

// WithBaseURL overrides the base URL.
func WithBaseURL(url string) Option {
	return func(options *Options) {
		options.baseURL = url
	}
}

// WithAPIKey enables API key authentication.
func WithAPIKey(apiKey string) Option {
	return func(options *Options) {
		options.auth = &authentication{
			token:  strings.TrimSpace(apiKey),
			method: apiKeyAuthentication,
		}
	}
}

// WithOAuthToken enables Oauth token authentication.
func WithOAuthToken(token string) Option {
	return func(options *Options) {
		options.auth = &authentication{
			token:  strings.TrimSpace(token),
			method: oAuthAuthentication,
		}
	}
}
