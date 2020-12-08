package internal

import (
	"fmt"
	"net/mail"
	"strings"

	"github.com/araddon/dateparse"

	"github.com/xitonix/createsend/transactional"
)

// SmartEmailDetails raw smart email details model.
type SmartEmailDetails struct {
	// ID email ID.
	ID string `json:"SmartEmailID"`
	// Name email name.
	Name string
	// CreatedAt the time when the email was created.
	CreatedAt string
	// Status the current status of the email.
	Status transactional.SmartEmailStatus
	// Properties smart email properties
	Properties struct {
		// From sender's email address
		From string
		// ReplyTo email address
		ReplyTo string
		// Subject subject line
		Subject string
		// Content content properties
		Content struct {
			// HTML HTML content
			HTML string
			// Text Text content
			Text string
			// EmailVariables email variables
			EmailVariables []string
			// InlineCSS inline CSS
			InlineCSS bool
		}
		// TextPreviewURL the URL to preview the text version
		TextPreviewURL string
		// HTMLPreviewURL the URL to preview the HTML version
		HTMLPreviewURL string
	}
	// AddRecipientsToList the optional ID of a subscriber list to which all recipients will be added
	AddRecipientsToList string
}

// ToSmartEmailDetails converts the raw model to a new createsend model.
func (s *SmartEmailDetails) ToSmartEmailDetails() (*transactional.SmartEmailDetails, error) {
	if s == nil || s.ID == "" {
		return &transactional.SmartEmailDetails{}, nil
	}

	date, err := dateparse.ParseAny(s.CreatedAt)
	if err != nil {
		return nil, err
	}

	from, err := mail.ParseAddress(s.Properties.From)
	if err != nil {
		return nil, fmt.Errorf("invalid From address: %w", err)
	}

	var replyTo *mail.Address
	if len(strings.TrimSpace(s.Properties.ReplyTo)) > 0 {
		replyTo, err = mail.ParseAddress(s.Properties.ReplyTo)
		if err != nil {
			return nil, fmt.Errorf("invalid ReplyTo address: %w", err)
		}
	}

	return &transactional.SmartEmailDetails{
		SmartEmailBasicDetails: transactional.SmartEmailBasicDetails{
			ID:        s.ID,
			Name:      s.Name,
			CreatedAt: date,
			Status:    s.Status,
		},
		From:                *from,
		ReplyTo:             replyTo,
		Subject:             s.Properties.Subject,
		HTML:                s.Properties.Content.HTML,
		Text:                s.Properties.Content.Text,
		EmailVariables:      s.Properties.Content.EmailVariables,
		InlineCSS:           s.Properties.Content.InlineCSS,
		TextPreviewURL:      s.Properties.TextPreviewURL,
		HTMLPreviewURL:      s.Properties.HTMLPreviewURL,
		AddRecipientsToList: s.AddRecipientsToList,
	}, nil
}
