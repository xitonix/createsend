package createsend

import (
	"github.com/xitonix/createsend/accounts"
	"github.com/xitonix/createsend/campaigns"
	"github.com/xitonix/createsend/clients"
	"github.com/xitonix/createsend/transactional"
)

// Client represents a client to access Campaign Monitor API.
type Client struct {
	accounts      accounts.API
	clients       clients.API
	transactional transactional.API
	campaigns     campaigns.API
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

	if opts.transactional == nil {
		client.transactional = newTransactionalAPI(hc)
	}

	if opts.campaigns == nil {
		client.campaigns = newCampaignAPI(hc)
	}

	return client, nil
}

// Accounts accesses Campaign Monitor's accounts API.
func (c *Client) Accounts() accounts.API {
	return c.accounts
}

// Clients accesses Campaign Monitor's clients API.
func (c *Client) Clients() clients.API {
	return c.clients
}

// Transactional accesses Campaign Monitor's Transactional API.
func (c *Client) Transactional() transactional.API {
	return c.transactional
}

// Campaigns accesses the Campaign Monitor campaigns API.
func (c *Client) Campaigns() campaigns.API {
	return c.campaigns
}
