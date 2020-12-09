package transactional

import (
	"net/mail"
	"time"
)

// SmartEmailBasicDetails represents the basic details of a Smart Transactional email.
type SmartEmailBasicDetails struct {
	// ID email ID.
	ID string
	// Name email name.
	Name string
	// CreatedAt the time when the email was created.
	CreatedAt time.Time
	// Status the current status of the email.
	Status SmartEmailStatus
}

// SmartEmailDetails represents the details of a Smart Transactional email.
type SmartEmailDetails struct {
	SmartEmailBasicDetails
	// From sender's email address
	From mail.Address
	// ReplyTo email address
	ReplyTo *mail.Address
	// Subject subject line
	Subject string
	// HTML HTML content
	HTML string
	// Text Text content
	Text string
	// EmailVariables email variables
	EmailVariables []string
	// InlineCSS inline CSS
	InlineCSS bool
	// TextPreviewURL the URL to preview the text version
	TextPreviewURL string
	// HTMLPreviewURL the URL to preview the HTML version
	HTMLPreviewURL string
	// AddRecipientsToList the optional ID of a subscriber list to which all recipients will be added
	AddRecipientsToList string
}
