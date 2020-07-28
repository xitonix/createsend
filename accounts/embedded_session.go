package accounts

// EmbeddedSession represents a login session for the member with the specified email address.
//
// See https://www.campaignmonitor.com/api/account/#single-sign-on for more details.
type EmbeddedSession struct {
	// EmailAddress a valid email address for a person's account in the selected Campaign Monitor client.
	EmailAddress string `json:"Email"`
	// Chrome defines what Campaign Monitor navigation to show. Valid options are "All", "Tabs", or "None".
	Chrome string
	// URL the Campaign Monitor page to load
	URL string `json:"Url"`
	// IntegratorID your integration ID.
	IntegratorID string
	// ClientID the client ID of the account you want to access.
	ClientID string
}
