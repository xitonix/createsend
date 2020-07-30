package campaigns

// Draft represents a draft campaign.
type Draft struct {
	// Name campaign name.
	Name string
	// Subject campaign subject.
	Subject string
	// FromName sender's name.
	FromName string
	// FromEmail sender's email address.
	FromEmail string
	// ReplyTo the reply-to email address.
	ReplyTo string
	// Text the optional URL of the text content.
	Text string `json:"TextUrl,omitempty"`
	// HTML the URL of the HTML content (eg. http://example.com/campaigncontent/index.html)
	HTML string `json:"HtmlUrl" `
	// Lists the ID of the lists you’d like the campaign to be eventually sent to.
	Lists []string `json:"ListIDs,omitempty"`
	// Segments the ID of the segments you’d like the campaign to be eventually sent to.
	Segments []string `json:"SegmentIDs,omitempty"`
}
