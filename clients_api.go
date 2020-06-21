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

	result := new(struct {
		ApiKey       string
		BasicDetails struct {
			ClientID     string
			CompanyName  string
			Country      string
			TimeZone     string
			EmailAddress string
			ContactName  string
		}
		AccessDetails *struct {
			AccessLevel int
			Username    string
		}
		BillingDetails        *internal.BillingDetails
		PendingBillingDetails *internal.BillingDetails
	})
	err := a.client.Get(path, &result)
	if err != nil {
		return nil, err
	}

	clientDetails := &clients.ClientDetails{
		APIKey:   result.ApiKey,
		Id:       result.BasicDetails.ClientID,
		Company:  result.BasicDetails.CompanyName,
		Country:  result.BasicDetails.Country,
		Timezone: result.BasicDetails.TimeZone,
		Billing:  result.BillingDetails.ToClientBillingDetails(result.PendingBillingDetails),
	}

	if result.BasicDetails.ContactName != "" || result.BasicDetails.EmailAddress != "" {
		clientDetails.Contact = &clients.ContactDetails{
			Name:         result.BasicDetails.ContactName,
			EmailAddress: result.BasicDetails.EmailAddress,
			AccessLevel:  -1,
			Username:     "",
		}
	}

	if result.AccessDetails != nil {
		if clientDetails.Contact == nil {
			clientDetails.Contact = &clients.ContactDetails{}
		}
		clientDetails.Contact.AccessLevel = result.AccessDetails.AccessLevel
		clientDetails.Contact.Username = result.AccessDetails.Username
	}

	return clientDetails, nil
}

func (a *clientsAPI) SentCampaigns(clientId string) ([]*clients.SentCampaign, error) {
	result := make([]*internal.SentCampaign, 0)
	path := fmt.Sprintf("clients/%s/campaigns.json", url.QueryEscape(clientId))
	err := a.client.Get(path, &result)
	if err != nil {
		return nil, err
	}
	campaigns := make([]*clients.SentCampaign, len(result))
	for i, c := range result {
		campaigns[i] = c.ToSendCampaign()
	}
	return campaigns, nil
}
