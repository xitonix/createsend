package createsend

import (
	"errors"
	"fmt"
	"github.com/xitonix/createsend/campaigns"
	"github.com/xitonix/createsend/campaigns/orderfield"
	"github.com/xitonix/createsend/internal"
	"github.com/xitonix/createsend/order"
	"net/url"
	"time"
)

type campaignAPI struct {
	client internal.Client
}

func newCampaignAPI(client internal.Client) *campaignAPI {
	return &campaignAPI{client: client}
}

func (c *campaignAPI) Create(clientID string, campaign campaigns.CampaignFromUrl) (string, error) {
	path := fmt.Sprintf("campaigns/%s.json", url.QueryEscape(clientID))
	var cId string
	err := c.client.Post(path, &cId, campaign)
	return cId, err
}

func (c *campaignAPI) CreateFromTemplate(clientID string, campaign campaigns.CampaignFromTemplate) (string, error) {
	path := fmt.Sprintf("campaigns/%s/fromtemplate.json", url.QueryEscape(clientID))
	var cId string
	err := c.client.Post(path, &cId, campaign)
	return cId, err
}

func (c *campaignAPI) SendImmediately(campaignID string, confirmationEmails string) error {
	return Send(c.client, campaignID, confirmationEmails, "Immediately")
}

func (c *campaignAPI) ScheduleSend(campaignID string, confirmationEmails string, date time.Time) error {
	return Send(c.client, campaignID, confirmationEmails, date.Format("2006-01-02 15:04"))
}

func (c *campaignAPI) Test(campaignID string, previewRecipients []string) error {
	path := fmt.Sprintf("campaigns/%s/sendpreview.json", url.QueryEscape(campaignID))
	return c.client.Post(path, nil, struct {
		PreviewRecipients []string
	}{
		PreviewRecipients: previewRecipients,
	})
}

func (c *campaignAPI) Summary(campaignID string) (campaigns.CampaignSummary, error) {
	path := fmt.Sprintf("campaigns/%s/summary.json", url.QueryEscape(campaignID))
	var cs campaigns.CampaignSummary
	err := c.client.Get(path, &cs)
	return cs, err
}

func (c *campaignAPI) EmailClientUsage(campaignID string) ([]campaigns.EmailClientUsage, error) {
	path := fmt.Sprintf("campaigns/%s/emailclientusage.json", url.QueryEscape(campaignID))
	var ecu []campaigns.EmailClientUsage
	err := c.client.Get(path, &ecu)
	return ecu, err
}

func (c *campaignAPI) ListsAndSegments(campaignID string) (campaigns.ListsAndSegments, error) {
	path := fmt.Sprintf("campaigns/%s/listsandsegments.json", url.QueryEscape(campaignID))
	var ls campaigns.ListsAndSegments
	err := c.client.Get(path, &ls)
	return ls, err
}

func (c *campaignAPI) Recipients(campaignID string, page int, pageSize int, orderField orderfield.OrderField, orderDirection order.Direction) (campaigns.Recipients, error) {
	var r campaigns.Recipients
	if orderField == orderfield.Date {
		return r, errors.New("date is not a valid order field")
	}

	path := resultsPath("bounces", campaignID, time.Time{}, page, pageSize, orderField, orderDirection)
	err := c.client.Get(path, &r)
	return r, err
}

func (c *campaignAPI) Bounces(campaignID string, date time.Time, page int, pageSize int, orderField orderfield.OrderField, orderDirection order.Direction) (campaigns.Bounces, error) {
	path := resultsPath("bounces", campaignID, date, page, pageSize, orderField, orderDirection)
	var b campaigns.Bounces
	err := c.client.Get(path, &b)
	return b, err
}

func (c *campaignAPI) Opens(campaignID string, date time.Time, page int, pageSize int, orderField orderfield.OrderField, orderDirection order.Direction) (campaigns.CampaignRecipientActions, error) {
	path := resultsPath("opens", campaignID, date, page, pageSize, orderField, orderDirection)
	var o campaigns.CampaignRecipientActions
	err := c.client.Get(path, &o)
	return o, err
}

func (c *campaignAPI) Clicks(campaignID string, date time.Time, page int, pageSize int, orderField orderfield.OrderField, orderDirection order.Direction) (campaigns.CampaignRecipientActions, error) {
	path := resultsPath("clicks", campaignID, date, page, pageSize, orderField, orderDirection)
	var o campaigns.CampaignRecipientActions
	err := c.client.Get(path, &o)
	return o, err
}

func (c *campaignAPI) Unsubscribes(campaignID string, date time.Time, page int, pageSize int, orderField orderfield.OrderField, orderDirection order.Direction) (campaigns.Unsubscribes, error) {
	path := resultsPath("unsubscribes", campaignID, date, page, pageSize, orderField, orderDirection)
	var u campaigns.Unsubscribes
	err := c.client.Get(path, &u)
	return u, err
}

func (c *campaignAPI) SpamComplaints(campaignID string, date time.Time, page int, pageSize int, orderField orderfield.OrderField, orderDirection order.Direction) (campaigns.SpamComplaints, error) {
	path := resultsPath("spam", campaignID, date, page, pageSize, orderField, orderDirection)
	var s campaigns.SpamComplaints
	err := c.client.Get(path, &s)
	return s, err
}

func (c *campaignAPI) Delete(campaignID string) error {
	return c.client.Delete(fmt.Sprintf("campaigns/%s.json", campaignID))
}

func (c *campaignAPI) Unschedule(campaignID string) error {
	return c.client.Delete(fmt.Sprintf("campaigns/%s/unschedule.json", campaignID))
}

func resultsPath(action string, campaignID string, date time.Time, page int, pageSize int, orderField orderfield.OrderField, orderDirection order.Direction) string {
	var dateQueryString = ""
	if !date.IsZero() {
		dateQueryString = fmt.Sprintf("date=%s", date.Format("2006-01-02 15:04"))
	}

	return fmt.Sprintf("campaigns/%s/%s.json?%spage=%d&pagesize=%d&orderfield=%s&orderdirection=%s", url.QueryEscape(campaignID), action, dateQueryString, page, pageSize, orderField, orderDirection)
}

func Send(client internal.Client, campaignID string, confirmationEmails string, formattedDate string) error {
	path := fmt.Sprintf("campaigns/%s/send.json", url.QueryEscape(campaignID))
	return client.Post(path, nil, struct {
		ConfirmationEmail string
		SendDate          string
	}{
		ConfirmationEmail: confirmationEmails,
		SendDate:          formattedDate,
	})
}
