package createsend

import (
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/xitonix/createsend/campaigns"
	"github.com/xitonix/createsend/internal"
	"github.com/xitonix/createsend/order"
	"net/url"
	"strings"
	"time"
)

type campaignsAPI struct {
	client internal.Client
}

func newCampaignAPI(client internal.Client) *campaignsAPI {
	return &campaignsAPI{client: client}
}

func (c *campaignsAPI) Create(clientID string, campaign campaigns.WithURLs) (string, error) {
	path := fmt.Sprintf("campaigns/%s.json", url.QueryEscape(clientID))
	var cId string
	err := c.client.Post(path, &cId, campaign)
	if err != nil {
		return "", err
	}

	return cId, nil
}

func (c *campaignsAPI) CreateFromTemplate(clientID string, campaign campaigns.Template) (string, error) {
	path := fmt.Sprintf("campaigns/%s/fromtemplate.json", url.QueryEscape(clientID))
	var cId string
	err := c.client.Post(path, &cId, campaign)
	if err != nil {
		return "", err
	}
	return cId, nil
}

func (c *campaignsAPI) Send(draftCampaignID string, confirmationEmails ...string) error {
	return c.SendAt(draftCampaignID, time.Time{}, confirmationEmails...)
}

func (c *campaignsAPI) SendAt(draftCampaignID string, at time.Time, confirmationEmails ...string) error {
	request := struct {
		ConfirmationEmail string
		SendDate          string
	}{
		ConfirmationEmail: strings.Join(confirmationEmails, ","),
		SendDate:          "Immediately",
	}
	if !at.IsZero() {
		request.SendDate = at.Format("2006-01-02 15:04")
	}
	path := fmt.Sprintf("campaigns/%s/send.json", url.QueryEscape(draftCampaignID))

	return c.client.Post(path, nil, request)
}

func (c *campaignsAPI) SendPreview(draftCampaignID string, recipients ...string) error {
	request := struct {
		PreviewRecipients []string
	}{
		PreviewRecipients: recipients,
	}
	path := fmt.Sprintf("campaigns/%s/sendpreview.json", url.QueryEscape(draftCampaignID))

	return c.client.Post(path, nil, request)
}

func (c *campaignsAPI) Summary(campaignID string) (campaigns.Summary, error) {
	path := fmt.Sprintf("campaigns/%s/summary.json", url.QueryEscape(campaignID))
	var cs campaigns.Summary
	err := c.client.Get(path, &cs)
	if err != nil {
		return campaigns.Summary{}, err
	}

	return cs, nil
}

func (c *campaignsAPI) EmailClientUsage(campaignID string) ([]campaigns.EmailClientUsage, error) {
	path := fmt.Sprintf("campaigns/%s/emailclientusage.json", url.QueryEscape(campaignID))
	var ecu []campaigns.EmailClientUsage
	err := c.client.Get(path, &ecu)
	if err != nil {
		return nil, err
	}

	return ecu, nil
}

func (c *campaignsAPI) ListsAndSegments(campaignID string) (campaigns.ListsAndSegments, error) {
	path := fmt.Sprintf("campaigns/%s/listsandsegments.json", url.QueryEscape(campaignID))
	var ls campaigns.ListsAndSegments
	err := c.client.Get(path, &ls)
	if err != nil {
		return campaigns.ListsAndSegments{}, err
	}

	return ls, nil
}

func (c *campaignsAPI) Recipients(campaignID string, page int, pageSize int, orderField order.Field, orderDirection order.Direction) (campaigns.Recipients, error) {
	if orderField == order.Date {
		return campaigns.Recipients{}, newClientError(ErrCodeInvalidDateOrderField)
	}

	var t struct {
		Results []struct {
			campaigns.Recipient
		}
		ResultsOrderedBy order.Field
		order.Page
	}

	path := getRecipientActivityPath("recipients", campaignID, time.Time{}, page, pageSize, orderField, orderDirection)
	err := c.client.Get(path, &t)
	if err != nil {
		return campaigns.Recipients{}, err
	}

	r := campaigns.Recipients{
		Results:   make([]campaigns.Recipient, len(t.Results)),
		OrderedBy: t.ResultsOrderedBy,
		Page:      t.Page,
	}
	for i := 0; i < len(t.Results); i++ {
		r.Results[i] = t.Results[i].Recipient
	}

	return r, nil
}

func (c *campaignsAPI) Bounces(campaignID string, date time.Time, page int, pageSize int, orderField order.Field, orderDirection order.Direction) (campaigns.Bounces, error) {
	path := getRecipientActivityPath("bounces", campaignID, date, page, pageSize, orderField, orderDirection)
	var b campaigns.Bounces
	err := c.client.Get(path, &b)
	if err != nil {
		return campaigns.Bounces{}, err
	}

	return b, nil
}

func (c *campaignsAPI) Opens(campaignID string, date time.Time, page int, pageSize int, orderField order.Field, orderDirection order.Direction) (campaigns.Opens, error) {
	path := getRecipientActivityPath("opens", campaignID, date, page, pageSize, orderField, orderDirection)
	var t struct {
		Results []struct {
			campaigns.Recipient
			campaigns.OpenDetails
			Date string
		}
		ResultsOrderedBy order.Field
		order.Page
	}
	err := c.client.Get(path, &t)
	if err != nil {
		return campaigns.Opens{}, err
	}

	op := campaigns.Opens{
		Results:   make([]campaigns.OpenDetails, len(t.Results)),
		OrderedBy: t.ResultsOrderedBy,
		Page:      t.Page,
	}
	for i := 0; i < len(t.Results); i++ {
		op.Results[i].Recipient = t.Results[i].Recipient
		op.Results[i].Date, err = dateparse.ParseAny(t.Results[i].Date)
		op.Results[i].IPAddress = t.Results[i].IPAddress
		op.Results[i].Latitude = t.Results[i].Latitude
		op.Results[i].Longitude = t.Results[i].Longitude
		op.Results[i].City = t.Results[i].City
		op.Results[i].Region = t.Results[i].Region
		op.Results[i].CountryCode = t.Results[i].CountryCode
		op.Results[i].CountryName = t.Results[i].CountryName
		if err != nil {
			return campaigns.Opens{}, err
		}
	}

	return op, nil
}

func (c *campaignsAPI) Clicks(campaignID string, date time.Time, page int, pageSize int, orderField order.Field, orderDirection order.Direction) (campaigns.Clicks, error) {
	path := getRecipientActivityPath("clicks", campaignID, date, page, pageSize, orderField, orderDirection)
	var t struct {
		Results []struct {
			campaigns.Recipient
			campaigns.ClickDetails
			Date string
			URL  string
		}
		ResultsOrderedBy order.Field
		order.Page
	}
	err := c.client.Get(path, &t)
	if err != nil {
		return campaigns.Clicks{}, err
	}

	cl := campaigns.Clicks{
		Results:   make([]campaigns.ClickDetails, len(t.Results)),
		OrderedBy: t.ResultsOrderedBy,
		Page:      t.Page,
	}
	for i := 0; i < len(t.Results); i++ {
		cl.Results[i].Recipient = t.Results[i].Recipient
		cl.Results[i].Date, err = dateparse.ParseAny(t.Results[i].Date)
		cl.Results[i].URL = t.Results[i].URL
		cl.Results[i].IPAddress = t.Results[i].IPAddress
		cl.Results[i].Latitude = t.Results[i].Latitude
		cl.Results[i].Longitude = t.Results[i].Longitude
		cl.Results[i].City = t.Results[i].City
		cl.Results[i].Region = t.Results[i].Region
		cl.Results[i].CountryCode = t.Results[i].CountryCode
		cl.Results[i].CountryName = t.Results[i].CountryName
		if err != nil {
			return campaigns.Clicks{}, err
		}
	}

	return cl, nil
}

func (c *campaignsAPI) Unsubscribes(campaignID string, date time.Time, page int, pageSize int, orderField order.Field, orderDirection order.Direction) (campaigns.Unsubscribes, error) {
	path := getRecipientActivityPath("unsubscribes", campaignID, date, page, pageSize, orderField, orderDirection)
	var t struct {
		Results []struct {
			campaigns.Recipient
			Date      string
			IPAddress string
		}
		ResultsOrderedBy order.Field
		order.Page
	}
	err := c.client.Get(path, &t)
	if err != nil {
		return campaigns.Unsubscribes{}, err
	}

	u := campaigns.Unsubscribes{
		Results:   make([]campaigns.Unsubscribe, len(t.Results)),
		OrderedBy: t.ResultsOrderedBy,
		Page:      t.Page,
	}
	for i := 0; i < len(t.Results); i++ {
		u.Results[i].Recipient = t.Results[i].Recipient
		u.Results[i].IPAddress = t.Results[i].IPAddress
		u.Results[i].Date, err = dateparse.ParseAny(t.Results[i].Date)
		if err != nil {
			return campaigns.Unsubscribes{}, err
		}
	}

	return u, nil
}

func (c *campaignsAPI) SpamComplaints(campaignID string, date time.Time, page int, pageSize int, orderField order.Field, orderDirection order.Direction) (campaigns.SpamComplaints, error) {
	path := getRecipientActivityPath("spam", campaignID, date, page, pageSize, orderField, orderDirection)
	var t struct {
		Results []struct {
			campaigns.Recipient
			Date string
		}
		ResultsOrderedBy order.Field
		order.Page
	}
	err := c.client.Get(path, &t)
	if err != nil {
		return campaigns.SpamComplaints{}, err
	}

	s := campaigns.SpamComplaints{
		Results:   make([]campaigns.SpamComplaint, len(t.Results)),
		OrderedBy: t.ResultsOrderedBy,
		Page:      t.Page,
	}
	for i := 0; i < len(t.Results); i++ {
		s.Results[i].Recipient = t.Results[i].Recipient
		s.Results[i].Date, err = dateparse.ParseAny(t.Results[i].Date)
		if err != nil {
			return campaigns.SpamComplaints{}, err
		}
	}

	return s, nil
}

func (c *campaignsAPI) Delete(campaignID string) error {
	return c.client.Delete(fmt.Sprintf("campaigns/%s.json", campaignID))
}

func (c *campaignsAPI) Unschedule(campaignID string) error {
	return c.client.Delete(fmt.Sprintf("campaigns/%s/unschedule.json", campaignID))
}

func getRecipientActivityPath(action string, campaignID string, date time.Time, page int, pageSize int, orderField order.Field, orderDirection order.Direction) string {
	var dateQueryString = ""
	if !date.IsZero() {
		dateQueryString = fmt.Sprintf("date=%s", date.Format("2006-01-02 15:04"))
	}

	return fmt.Sprintf("campaigns/%s/%s.json?%spage=%d&pagesize=%d&orderfield=%s&orderdirection=%s", url.QueryEscape(campaignID), action, dateQueryString, page, pageSize, orderField, orderDirection)
}
