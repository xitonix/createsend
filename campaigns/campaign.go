package campaigns

// BasicDetails the basic details of a campaign
type BasicDetails struct {
	// Name represents the name of the campaign
	Name string
	// Subject represents the subject of the campaign
	Subject string
	// FromName represents the from name of the campaign
	FromName string
	// FromEmail represents the from email address
	FromEmail string
	// ReplyTo represents the reply to email address
	ReplyTo string
	// ListIds represents the list ids of the campaign
	ListIds []string
	// SegmentIds represents the segment ids of the campaign
	SegmentIds []string
}

// WithURLs the definition of a campaign
type WithURLs struct {
	// BasicDetails the basic details of a campaign
	BasicDetails
	// Html the URL of the HTML content
	Html string `json:"HtmlUrl"`
	// Text the URL of the text content
	Text string `json:"TextUrl"`
}
