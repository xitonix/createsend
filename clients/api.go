package clients

// API is an interface that wraps client related operations.
//
// The API contains all the functionality you need to manage the clients in your account.
// From adding new clients, updating their billing settings, giving them access to their account, accessing their lists,
// templates and campaigns, etc.
type API interface {
	// Create creates a new client in your account with basic contact information and no access to the application.
	//
	// Client billing options are set once the client is created.
	Create(client Client) (string, error)
	// Get returns the complete details for a client including their API key, access level, contact details and billing settings.
	Get(clientId string) (*ClientDetails, error)
	// SentCampaign returns a list of all sent campaigns for a client.
	SentCampaigns(clientId string) ([]*SentCampaign, error)
	// ScheduledCampaigns returns all currently scheduled campaigns for a client.
	ScheduledCampaigns(clientId string) ([]*ScheduledCampaign, error)
	// DraftCampaigns returns all draft campaigns belonging to a client.
	DraftCampaigns(clientId string) ([]*DraftCampaign, error)
}
