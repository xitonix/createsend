package internal

import (
	"github.com/xitonix/createsend/common"
	"net/mail"

	"github.com/araddon/dateparse"
)

// Campaign represents a raw Campaign.
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

// SentCampaign represent a raw sent campaign.
type SentCampaign struct {
	Campaign
	SentDate        string
	TotalRecipients int64
}

// ToSendCampaign converts the raw model to a new createsend model.
func (c *SentCampaign) ToSendCampaign() (*common.SentCampaign, error) {
	if c == nil {
		return nil, nil
	}
	date, err := dateparse.ParseAny(c.SentDate)
	if err != nil {
		return nil, err
	}
	return &common.SentCampaign{
		Campaign: common.Campaign{
			ID:   c.CampaignID,
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

// ScheduledCampaign represents a raw scheduled campaign.
type ScheduledCampaign struct {
	Campaign
	// DateCreated the timestamp when the Campaign was created.
	DateCreated string
	// DateScheduled the timestamp when the Campaign will be sent.
	DateScheduled string
	// ScheduledTimeZone schedule timezone.
	ScheduledTimeZone string
}

// ToScheduledCampaign converts the raw model to a new createsend model.
func (c *ScheduledCampaign) ToScheduledCampaign() (*common.ScheduledCampaign, error) {
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
	return &common.ScheduledCampaign{
		Campaign: common.Campaign{
			ID:   c.CampaignID,
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

// DraftCampaign represents a raw draft campaign.
type DraftCampaign struct {
	Campaign
	// DateCreated the timestamp when the Campaign was created.
	DateCreated string
}

// ToDraftCampaign converts the raw model to a new createsend model.
func (c *DraftCampaign) ToDraftCampaign() (*common.DraftCampaign, error) {
	date, err := dateparse.ParseAny(c.DateCreated)
	if err != nil {
		return nil, err
	}
	return &common.DraftCampaign{
		Campaign: common.Campaign{
			ID:   c.CampaignID,
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
