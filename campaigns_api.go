package createsend

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/xitonix/createsend/campaigns"
	"github.com/xitonix/createsend/internal"
)

type campaignsAPI struct {
	client internal.Client
}

func newCampaignsAPI(client internal.Client) *campaignsAPI {
	return &campaignsAPI{client: client}
}

func (a *campaignsAPI) CreateDraft(clientID string, campaign campaigns.Draft) (string, error) {
	path := fmt.Sprintf("campaigns/%s.json", url.QueryEscape(clientID))
	var result string
	err := a.client.Post(path, &result, campaign)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (a *campaignsAPI) Send(draftCampaignID string, confirmationEmails ...string) error {
	return a.SendAt(draftCampaignID, time.Time{}, confirmationEmails...)
}

func (a *campaignsAPI) SendAt(draftCampaignID string, at time.Time, confirmationEmails ...string) error {
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
	return a.client.Post(path, nil, request)
}

func (a *campaignsAPI) SendPreview(draftCampaignID string, recipients ...string) error {
	request := struct {
		PreviewRecipients []string
	}{
		PreviewRecipients: recipients,
	}
	path := fmt.Sprintf("campaigns/%s/sendpreview.json", url.QueryEscape(draftCampaignID))
	return a.client.Post(path, nil, request)
}
