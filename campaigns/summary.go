package campaigns

type Summary struct {
	Recipients        int
	TotalOpened       int
	Clicks            int
	Unsubscribed      int
	Bounced           int
	UniqueOpened      int
	SpamComplaints    int
	WebVersionURL     string
	WebVersionTextURL string
	WorldviewURL      string
	Forwards          int
	Likes             int
	Mentions          int
}
