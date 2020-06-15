package createsend

import (
	"fmt"
	"net/url"
	"time"

	"github.com/xitonix/createsend/accounts"
	"github.com/xitonix/createsend/internal"
)

const (
	listClientsPath         = "clients.json"
	fetchBillingDetailsPath = "billingdetails.json"
	fetchValidCountriesPath = "countries.json"
	fetchValidTimezonesPath = "timezones.json"
	fetchCurrentDatePath    = "systemdate.json"
	administratorsPath      = "admins.json"
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

func (a *accountsAPI) Now() (time.Time, error) {
	var result *struct {
		SystemDate string
	}
	err := a.client.Get(fetchCurrentDatePath, &result)
	if err != nil {
		return time.Time{}, err
	}
	if result != nil && len(result.SystemDate) > 0 {
		t, err := time.Parse("2006-01-02 15:04:05", result.SystemDate)
		if err != nil {
			return time.Time{}, newWrappedClientError("Failed to parse the server date value", err, ErrCodeDataProcessing)
		}
		return t, nil
	}

	return time.Time{}, nil
}

func (a *accountsAPI) AddAdministrator(administrator accounts.Administrator) error {
	return a.client.Post(administratorsPath, nil, administrator)
}

func (a *accountsAPI) UpdateAdministrator(currentEmailAddress string, administrator accounts.Administrator) error {
	path := fmt.Sprintf("%s?email=%s", administratorsPath, url.QueryEscape(currentEmailAddress))
	return a.client.Put(path, nil, administrator)
}

func (a *accountsAPI) Administrators() ([]*accounts.AdministratorDetails, error) {
	result := make([]*accounts.AdministratorDetails, 0)
	err := a.client.Get(administratorsPath, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *accountsAPI) Administrator(emailAddress string) (*accounts.AdministratorDetails, error) {
	var result *accounts.AdministratorDetails
	path := fmt.Sprintf("%s?email=%s", administratorsPath, url.QueryEscape(emailAddress))
	err := a.client.Get(path, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
