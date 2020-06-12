package createsend

import (
	"github.com/xitonix/createsend/accounts"
	"github.com/xitonix/createsend/internal"
)

type accountsAPI struct {
	client internal.Client
}

func newAccountAPI(client internal.Client) *accountsAPI {
	return &accountsAPI{client: client}
}

func (a *accountsAPI) Clients() ([]*accounts.Client, error) {
	result := make([]*accounts.Client, 0)
	err := a.client.Get("clients.json", &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *accountsAPI) Billing() (*accounts.Billing, error) {
	var result *accounts.Billing
	err := a.client.Get("billingdetails.json", &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
