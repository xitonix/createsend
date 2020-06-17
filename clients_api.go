package createsend

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/xitonix/createsend/clients"
	"github.com/xitonix/createsend/internal"
)

const (
	clientsPath = "clients.json"
)

type clientsAPI struct {
	client internal.Client
}

func newClientsAPI(client internal.Client) *clientsAPI {
	return &clientsAPI{client: client}
}

func (a *clientsAPI) Create(client clients.Client) (string, error) {
	var clientId string
	err := a.client.Post(clientsPath, &clientId, client)
	return strings.Trim(clientId, `"`), err
}

func (a *clientsAPI) Get(clientId string) (*clients.ClientDetails, error) {
	path := fmt.Sprintf("clients/%s.json", url.QueryEscape(clientId))
	var result *struct {
		ApiKey       string
		BasicDetails struct {
			ClientID            string
			CompanyName         string
			Country             string
			TimeZone            string
			PrimaryContactName  string
			PrimaryContactEmail string
		}
		BillingDetails clients.BillingDetails
	}
	err := a.client.Get(path, &result)
	if err != nil {
		return nil, err
	}
	return &clients.ClientDetails{
		ApiKey:              result.ApiKey,
		Id:                  result.BasicDetails.ClientID,
		Company:             result.BasicDetails.CompanyName,
		Country:             result.BasicDetails.Country,
		Timezone:            result.BasicDetails.TimeZone,
		PrimaryContactName:  result.BasicDetails.PrimaryContactName,
		PrimaryContactEmail: result.BasicDetails.PrimaryContactEmail,
		Billing:             result.BillingDetails,
	}, nil
}
