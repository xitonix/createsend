package clients

import (
	"net/mail"
)

// Campaign represents a Campaign.
type Campaign struct {
	// Id campaign id.
	Id string
	// Name campaign name.
	Name string
	// From the sender.
	From mail.Address
	// ReplyTo reply to value.
	ReplyTo string
	// WebVersionURL web version URL.
	WebVersionURL string
	// WebVersionTextURL the plain text format of the web version URL.
	WebVersionTextURL string
	// Subject Campaign's subject.
	Subject string
}

// SentCampaign represents a sent Campaign.
type SentCampaign struct {
	// Campaign Campaign's basic details.
	Campaign
	// SentDate the timestamp when the Campaign was sent.
	SentDate string
	// Recipients number of recipients the Campaign was sent to.
	Recipients int64
}

// ScheduledCampaign represents a scheduled Campaign.
type ScheduledCampaign struct {
	// Campaign Campaign's basic details.
	Campaign
	// DateCreated the timestamp when the Campaign was created.
	DateCreated string
	// DateScheduled the timestamp when the Campaign will be sent.
	DateScheduled string
	// Timezone schedule timezone.
	Timezone string
}
