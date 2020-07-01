package internal

import (
	"net/mail"

	"github.com/araddon/dateparse"

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
}

type SentCampaign struct {
	Campaign
	SentDate        string
	TotalRecipients int64
}

func (c *SentCampaign) ToSendCampaign() (*clients.SentCampaign, error) {
	if c == nil {
		return nil, nil
	}
	date, err := dateparse.ParseAny(c.SentDate)
	if err != nil {
		return nil, err
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
		SentDate:   date,
		Recipients: c.TotalRecipients,
	}, nil
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

func (c *ScheduledCampaign) ToScheduledCampaign() (*clients.ScheduledCampaign, error) {
	if c == nil {
		return nil, nil
	}
	dc, err := dateparse.ParseAny(c.DateCreated)
	if err != nil {
		return nil, err
	}

	ds, err := dateparse.ParseAny(c.DateScheduled)
	if err != nil {
		return nil, err
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
		DateCreated:   dc,
		DateScheduled: ds,
		Timezone:      c.ScheduledTimeZone,
	}, nil
}

type DraftCampaign struct {
	Campaign
	// DateCreated the timestamp when the Campaign was created.
	DateCreated string
}

func (c *DraftCampaign) ToDraftCampaign() (*clients.DraftCampaign, error) {
	date, err := dateparse.ParseAny(c.DateCreated)
	if err != nil {
		return nil, err
	}
	return &clients.DraftCampaign{
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
		DateCreated: date,
	}, nil
}
