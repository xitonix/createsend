package createsend

import (
	"github.com/xitonix/createsend/accounts"
	"github.com/xitonix/createsend/campaigns"
	"github.com/xitonix/createsend/clients"
)

// Client represents a client to access Campaign Monitor API.
type Client struct {
	accounts  accounts.API
	clients   clients.API
	campaigns campaigns.API
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

	client := &Client{}
	if opts.accounts == nil {
		client.accounts = newAccountAPI(hc)
	}

	if opts.clients == nil {
		client.clients = newClientsAPI(hc)
	}

	if opts.campaigns == nil {
		client.campaigns = newCampaignAPI(hc)
	}

	return client, nil
}

// Accounts accesses the Campaign Monitor accounts API.
func (c *Client) Accounts() accounts.API {
	return c.accounts
}

// Clients accesses the Campaign Monitor clients API.
func (c *Client) Clients() clients.API {
	return c.clients
}

// Campaigns accesses the Campaign Monitor campaigns API.
func (c *Client) Campaigns() campaigns.API {
	return c.campaigns
}
