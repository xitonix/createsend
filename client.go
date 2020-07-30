package createsend

import (
	"github.com/xitonix/createsend/accounts"
	"github.com/xitonix/createsend/campaigns"
	"github.com/xitonix/createsend/clients"
)

// Client represents a client to access Campaign Monitor API.
type Client struct {
	opts *Options
}

// New creates a new client.
func New(options ...Option) (*Client, error) {
	opts := defaultOptions()
	for _, op := range options {
		op(opts)
	}

	hc, err := newHTTPClient(opts.ctx, opts.baseURL, opts.client, opts.auth)
	if err != nil {
		return nil, err
	}

	if opts.accounts == nil {
		opts.accounts = newAccountAPI(hc)
	}

	if opts.clients == nil {
		opts.clients = newClientsAPI(hc)
	}

	if opts.campaigns == nil {
		opts.campaigns = newCampaignsAPI(hc)
	}

	return &Client{
		opts: opts,
	}, nil
}

// Accounts accesses the Campaign Monitor Accounts API.
func (c *Client) Accounts() accounts.API {
	return c.opts.accounts
}

// Clients accesses the Campaign Monitor Clients API.
func (c *Client) Clients() clients.API {
	return c.opts.clients
}

// Campaigns accesses the Campaign Monitor Campaigns API.
func (c *Client) Campaigns() campaigns.API {
	return c.opts.campaigns
}
