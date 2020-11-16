package createsend

import (
	"fmt"
	"github.com/xitonix/createsend/common"
	"net/url"
	"strings"

	"github.com/xitonix/createsend/clients"
	"github.com/xitonix/createsend/internal"
	"github.com/xitonix/createsend/order"
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

func (a *clientsAPI) Create(details clients.BasicDetails) (string, error) {
	var clientID string
	err := a.client.Post(clientsPath, &clientID, details)
	return strings.Trim(clientID, `"`), err
}

func (a *clientsAPI) Get(clientID string) (*clients.ClientDetails, error) {
	path := fmt.Sprintf("clients/%s.json", url.QueryEscape(clientID))

	result := new(struct {
		APIKey       string
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
		APIKey:   result.APIKey,
		ID:       result.BasicDetails.ClientID,
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

func (a *clientsAPI) SentCampaigns(clientID string) ([]*common.SentCampaign, error) {
	result := make([]*internal.SentCampaign, 0)
	path := fmt.Sprintf("clients/%s/campaigns.json", url.QueryEscape(clientID))
	err := a.client.Get(path, &result)
	if err != nil {
		return nil, err
	}
	campaigns := make([]*common.SentCampaign, len(result))
	for i, c := range result {
		cm, err := c.ToSendCampaign()
		if err != nil {
			return nil, newClientError(ErrCodeDataProcessing)
		}
		campaigns[i] = cm
	}
	return campaigns, nil
}

func (a *clientsAPI) ScheduledCampaigns(clientID string) ([]*common.ScheduledCampaign, error) {
	result := make([]*internal.ScheduledCampaign, 0)
	path := fmt.Sprintf("clients/%s/scheduled.json", url.QueryEscape(clientID))
	err := a.client.Get(path, &result)
	if err != nil {
		return nil, err
	}
	campaigns := make([]*common.ScheduledCampaign, len(result))
	for i, c := range result {
		cm, err := c.ToScheduledCampaign()
		if err != nil {
			return nil, newClientError(ErrCodeDataProcessing)
		}
		campaigns[i] = cm
	}
	return campaigns, nil
}

func (a *clientsAPI) DraftCampaigns(clientID string) ([]*common.DraftCampaign, error) {
	result := make([]*internal.DraftCampaign, 0)
	path := fmt.Sprintf("clients/%s/drafts.json", url.QueryEscape(clientID))
	err := a.client.Get(path, &result)
	if err != nil {
		return nil, err
	}
	campaigns := make([]*common.DraftCampaign, len(result))
	for i, c := range result {
		cm, err := c.ToDraftCampaign()
		if err != nil {
			return nil, newClientError(ErrCodeDataProcessing)
		}
		campaigns[i] = cm
	}
	return campaigns, nil
}

func (a *clientsAPI) Lists(clientID string) ([]*clients.List, error) {
	result := make([]*clients.List, 0)
	path := fmt.Sprintf("clients/%s/lists.json", url.QueryEscape(clientID))
	err := a.client.Get(path, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *clientsAPI) ListsByEmailAddress(clientID, emailAddress string) ([]*clients.SubscriberList, error) {
	result := make([]*internal.SubscriberList, 0)
	path := fmt.Sprintf("clients/%s/listsforemail.json?email=%s", url.QueryEscape(clientID), url.QueryEscape(emailAddress))
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

func (a *clientsAPI) Segments(clientID string) ([]*clients.Segment, error) {
	result := make([]*clients.Segment, 0)
	path := fmt.Sprintf("clients/%s/segments.json", url.QueryEscape(clientID))
	err := a.client.Get(path, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *clientsAPI) SuppressionList(clientID string,
	pageSize, page int,
	orderBy order.SuppressionListField,
	direction order.Direction) (*clients.SuppressionList, error) {

	path := fmt.Sprintf("clients/%s/suppressionlist.json?page=%d&pagesize=%d&orderfield=%s&orderdirection=%s",
		url.QueryEscape(clientID),
		page,
		pageSize,
		url.QueryEscape(orderBy.String()),
		url.QueryEscape(direction.String()))

	result := new(internal.SuppressionList)
	err := a.client.Get(path, &result)
	if err != nil {
		return nil, err
	}

	list, err := result.ToSuppressionList()
	if err != nil {
		return nil, newClientError(ErrCodeDataProcessing)
	}

	return list, nil
}

func (a *clientsAPI) Suppress(clientID string, emails ...string) error {
	data := struct {
		EmailAddresses []string
	}{
		EmailAddresses: emails,
	}
	path := fmt.Sprintf("clients/%s/suppress.json", url.QueryEscape(clientID))
	return a.client.Post(path, nil, data)
}

func (a *clientsAPI) UnSuppress(clientID string, email string) error {
	path := fmt.Sprintf("clients/%s/unsuppress.json?email=%s",
		url.QueryEscape(clientID),
		url.QueryEscape(email))
	return a.client.Put(path, nil, nil)
}

func (a *clientsAPI) Templates(clientID string) ([]*clients.Template, error) {
	result := make([]*clients.Template, 0)
	path := fmt.Sprintf("clients/%s/templates.json", url.QueryEscape(clientID))
	err := a.client.Get(path, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *clientsAPI) Update(clientID string, details clients.BasicDetails) error {
	path := fmt.Sprintf("clients/%s/setbasics.json", url.QueryEscape(clientID))
	return a.client.Put(path, nil, details)
}

func (a *clientsAPI) SetPAYGBilling(clientID string, rates clients.PAYGRates) error {
	path := fmt.Sprintf("clients/%s/setpaygbilling.json", url.QueryEscape(clientID))
	return a.client.Put(path, nil, rates)
}

func (a *clientsAPI) SetMonthlyBilling(clientID string, rates clients.MonthlyRates) error {
	path := fmt.Sprintf("clients/%s/setmonthlybilling.json", url.QueryEscape(clientID))
	return a.client.Put(path, nil, rates)
}

func (a *clientsAPI) TransferCredits(clientID string, request clients.CreditTransferRequest) (*clients.CreditTransferResult, error) {
	path := fmt.Sprintf("clients/%s/credits.json", url.QueryEscape(clientID))
	var result *clients.CreditTransferResult
	err := a.client.Post(path, &result, request)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *clientsAPI) Delete(clientID string) error {
	path := fmt.Sprintf("clients/%s.json", url.QueryEscape(clientID))
	return a.client.Delete(path)
}

func (a *clientsAPI) AddPerson(clientID string, person clients.Person) (string, error) {
	path := fmt.Sprintf("clients/%s/people.json", url.QueryEscape(clientID))
	result := new(struct {
		EmailAddress string
	})
	err := a.client.Post(path, &result, person)
	if err != nil {
		return "", err
	}
	return result.EmailAddress, nil
}

func (a *clientsAPI) UpdatePerson(clientID string, emailAddress string, person clients.Person) (string, error) {
	path := fmt.Sprintf("clients/%s/people.json?email=%s", url.QueryEscape(clientID), url.QueryEscape(emailAddress))
	result := new(struct {
		EmailAddress string
	})
	err := a.client.Put(path, &result, person)
	if err != nil {
		return "", err
	}
	return result.EmailAddress, nil
}

func (a *clientsAPI) People(clientID string) ([]*clients.PersonDetails, error) {
	result := make([]*clients.PersonDetails, 0)
	path := fmt.Sprintf("clients/%s/people.json", url.QueryEscape(clientID))
	err := a.client.Get(path, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *clientsAPI) Person(clientID string, emailAddress string) (*clients.PersonDetails, error) {
	var result *clients.PersonDetails
	path := fmt.Sprintf("clients/%s/people.json?email=%s", url.QueryEscape(clientID), url.QueryEscape(emailAddress))
	err := a.client.Get(path, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *clientsAPI) DeletePerson(clientID string, emailAddress string) error {
	path := fmt.Sprintf("clients/%s/people.json?email=%s", url.QueryEscape(clientID), url.QueryEscape(emailAddress))
	return a.client.Delete(path)
}

func (a *clientsAPI) SetPrimaryContact(clientID string, emailAddress string) (string, error) {
	path := fmt.Sprintf("clients/%s/primarycontact.json?email=%s", url.QueryEscape(clientID), url.QueryEscape(emailAddress))
	result := new(struct {
		EmailAddress string
	})
	err := a.client.Put(path, &result, nil)
	if err != nil {
		return "", err
	}
	return result.EmailAddress, nil
}

func (a *clientsAPI) PrimaryContact(clientID string) (string, error) {
	path := fmt.Sprintf("clients/%s/primarycontact.json", url.QueryEscape(clientID))
	result := new(struct {
		EmailAddress string
	})
	err := a.client.Get(path, &result)
	if err != nil {
		return "", err
	}
	return result.EmailAddress, nil
}
