package campaigns

// Summary of a sent campaign
type Summary struct {
	// Recipients represents the amount of recipients of the campaign
	Recipients int
	// TotalOpened represents the total amount of emails that were opened
	TotalOpened int
	// UniqueOpened represents the total unique amount of email opens
	UniqueOpened int
	// Clicks represents the total amount of email clicks
	Clicks int
	// Clicks represents the total amount of email unsubscribes
	Unsubscribed int
	// Bounced represents the total amount of emails that bounced
	Bounced int
	// SpamComplaints represents the total amount of email spam complaints
	SpamComplaints int
	// WebVersionURL is the public URL of the campaign
	WebVersionURL string
	// WebVersionTextURL is the public URL of the text version of the campaign
	WebVersionTextURL string
	// WorldviewURL is a public URL linking to a world view showing details regarding the performance of the campaign
	WorldviewURL string
	// Forwards represents the total amount that the emails of the campaign was forwarded
	Forwards int
	// Likes represents the total amount of likes clicked on the campaign
	Likes int
	// Mentions represents the total amount of mentions clicked on the campaign
	Mentions int
}
