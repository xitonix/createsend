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
		cm, err := c.ToSendCampaign()
		if err != nil {
			return nil, newClientError(ErrCodeDataProcessing)
		}
		campaigns[i] = cm
	}
	return campaigns, nil
}

func (a *clientsAPI) ScheduledCampaigns(clientId string) ([]*clients.ScheduledCampaign, error) {
	result := make([]*internal.ScheduledCampaign, 0)
	path := fmt.Sprintf("clients/%s/scheduled.json", url.QueryEscape(clientId))
	err := a.client.Get(path, &result)
	if err != nil {
		return nil, err
	}
	campaigns := make([]*clients.ScheduledCampaign, len(result))
	for i, c := range result {
		cm, err := c.ToScheduledCampaign()
		if err != nil {
			return nil, newClientError(ErrCodeDataProcessing)
		}
		campaigns[i] = cm
	}
	return campaigns, nil
}

func (a *clientsAPI) DraftCampaigns(clientId string) ([]*clients.DraftCampaign, error) {
	result := make([]*internal.DraftCampaign, 0)
	path := fmt.Sprintf("clients/%s/drafts.json", url.QueryEscape(clientId))
	err := a.client.Get(path, &result)
	if err != nil {
		return nil, err
	}
	campaigns := make([]*clients.DraftCampaign, len(result))
	for i, c := range result {
		cm, err := c.ToDraftCampaign()
		if err != nil {
			return nil, newClientError(ErrCodeDataProcessing)
		}
		campaigns[i] = cm
	}
	return campaigns, nil
}

func (a *clientsAPI) Lists(clientId string) ([]*clients.List, error) {
	result := make([]*clients.List, 0)
	path := fmt.Sprintf("clients/%s/lists.json", url.QueryEscape(clientId))
	err := a.client.Get(path, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *clientsAPI) ListsByEmailAddress(clientId, emailAddress string) ([]*clients.SubscriberList, error) {
	result := make([]*internal.SubscriberList, 0)
	path := fmt.Sprintf("clients/%s/listsforemail.json?email=%s", url.QueryEscape(clientId), url.QueryEscape(emailAddress))
	err := a.client.Get(path, &result)
	if err != nil {
		return nil, err
	}
	lists := make([]*clients.SubscriberList, len(result))
	for i, r := range result {
		sl, err := r.ToSubscriberList()
		if err != nil {
			return nil, newClientError(ErrCodeDataProcessing)
		}
		lists[i] = sl
	}
	return lists, nil
}

func (a *clientsAPI) Segments(clientId string) ([]*clients.Segment, error) {
	result := make([]*clients.Segment, 0)
	path := fmt.Sprintf("clients/%s/segments.json", url.QueryEscape(clientId))
	err := a.client.Get(path, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
