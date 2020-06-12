package createsend

import (
	"github.com/xitonix/createsend/accounts"
)

type Client struct {
	accounts accounts.API
}

func New(options ...Option) (*Client, error) {
	opts := defaultOptions()
	for _, op := range options {
		op(opts)
	}

	hc, err := newHTTPClient(opts.baseURL, opts.client, opts.auth)
	if err != nil {
		return nil, err
	}

	client := &Client{}
	if opts.accounts == nil {
		client.accounts = newAccountAPI(hc)
	}

	return client, nil
}

func (c *Client) Accounts() accounts.API {
	return c.accounts
}
