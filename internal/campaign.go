package internal

import (
	"net/mail"

	"github.com/xitonix/createsend/clients"
)

type Campaign struct {
	FromName          string
	FromEmail         string
	ReplyTo           string
	WebVersionURL     string
	WebVersionTextURL string
	CampaignID        string
	Subject           string
	Name              string
	SentDate          string
	TotalRecipients   int64
}

type SentCampaign struct {
	Campaign
	SentDate        string
	TotalRecipients int64
}

func (c *SentCampaign) ToSendCampaign() *clients.SentCampaign {
	if c == nil {
		return nil
	}
	return &clients.SentCampaign{
		Campaign: clients.Campaign{
			Id:   c.CampaignID,
			Name: c.Name,
			From: mail.Address{
				Name:    c.FromName,
				Address: c.FromEmail,
			},
			ReplyTo:           c.ReplyTo,
			WebVersionURL:     c.WebVersionURL,
			WebVersionTextURL: c.WebVersionTextURL,
			Subject:           c.Subject,
		},
		SentDate:   c.SentDate,
		Recipients: c.TotalRecipients,
	}
}

type ScheduledCampaign struct {
	Campaign
	// DateCreated the timestamp when the Campaign was created.
	DateCreated string
	// DateScheduled the timestamp when the Campaign will be sent.
	DateScheduled string
	// ScheduledTimeZone schedule timezone.
	ScheduledTimeZone string
}

func (c *ScheduledCampaign) ToScheduledCampaign() *clients.ScheduledCampaign {
	if c == nil {
		return nil
	}
	return &clients.ScheduledCampaign{
		Campaign: clients.Campaign{
			Id:   c.CampaignID,
			Name: c.Name,
			From: mail.Address{
				Name:    c.FromName,
				Address: c.FromEmail,
			},
			ReplyTo:           c.ReplyTo,
			WebVersionURL:     c.WebVersionURL,
			WebVersionTextURL: c.WebVersionTextURL,
			Subject:           c.Subject,
		},
		DateCreated:   c.DateCreated,
		DateScheduled: c.DateScheduled,
		Timezone:      c.ScheduledTimeZone,
	}
}
