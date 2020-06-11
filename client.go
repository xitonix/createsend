package createsend

import (
	"github.com/gojektech/heimdall/v6/hystrix"
)

type Client struct {
	client *hystrix.Client
}

func New(options ...Option) (*Client, error) {
	opts := defaultOptions()
	for _, op := range options {
		op(opts)
	}

	httpClient, err := newHTTPClient(opts.baseURL, opts.client, opts.auth)
	if err != nil {
		return nil, err
	}

	client := hystrix.NewClient(
		hystrix.WithRetryCount(opts.retryCount),
		hystrix.WithHTTPClient(httpClient),
	)

	return &Client{client: client}, nil
}
