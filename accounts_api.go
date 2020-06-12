package createsend

import (
	"github.com/xitonix/createsend/accounts"
	"github.com/xitonix/createsend/internal"
)

const (
	listClientsPath         = "clients.json"
	fetchBillingDetailsPath = "billingdetails.json"
	fetchValidCountriesPath = "countries.json"
	fetchValidTimezonesPath = "timezones.json"
)

type accountsAPI struct {
	client internal.Client
}

func newAccountAPI(client internal.Client) *accountsAPI {
	return &accountsAPI{client: client}
}

func (a *accountsAPI) Clients() ([]*accounts.Client, error) {
	result := make([]*accounts.Client, 0)
	err := a.client.Get(listClientsPath, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *accountsAPI) Billing() (*accounts.Billing, error) {
	var result *accounts.Billing
	err := a.client.Get(fetchBillingDetailsPath, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *accountsAPI) Countries() ([]string, error) {
	var result []string
	err := a.client.Get(fetchValidCountriesPath, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *accountsAPI) Timezones() ([]string, error) {
	var result []string
	err := a.client.Get(fetchValidTimezonesPath, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
