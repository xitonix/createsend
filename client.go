package createsend

import (
	"github.com/xitonix/createsend/accounts"
	"github.com/xitonix/createsend/clients"
)

type Client struct {
	accounts accounts.API
	clients  clients.API
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

	if opts.clients == nil {
		client.clients = newClientsAPI(hc)
	}

	return client, nil
}

func (c *Client) Accounts() accounts.API {
	return c.accounts
}

func (c *Client) Clients() clients.API {
	return c.clients
}
