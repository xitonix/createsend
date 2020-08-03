package campaigns

// Summary represents a basic summary of the results for a sent campaign.
type Summary struct {
	// Recipients number of recipients the campaign was sent to.
	Recipients int
	// TotalOpened total number of recipients opened the campaign.
	TotalOpened int
	// UniqueOpened number of unique opens.
	UniqueOpened int
	// Clicks total number of clicks.
	Clicks int
	// Unsubscribed number of recipients unsubscribed from the campaign.
	Unsubscribed int
	// Bounced number of bounces.
	Bounced int
	// SpamComplaints number of recipients marked the campaign as spam.
	SpamComplaints int
	// WebVersionURL web version URL.
	WebVersionURL string
	// WebVersionTextURL text URL.
	WebVersionTextURL string
	// WorldviewURL world view URL.
	WorldviewURL string
	// Forwards number of times the campaign has been forwarded.
	Forwards int
	// Likes number of likes.
	Likes int
	// Mentions number of mentions.
	Mentions int
}
