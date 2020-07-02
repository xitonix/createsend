package clients

import "github.com/xitonix/createsend/order"

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
	// Lists returns all the subscriber lists that belong to a client.
	Lists(clientId string) ([]*List, error)
	// ListsByEmailAddress returns all the subscriber lists across the client, to which an email address is subscribed.
	ListsByEmailAddress(clientId, emailAddress string) ([]*SubscriberList, error)
	// Segments returns a list of all list segments belonging to a particular client.
	Segments(clientId string) ([]*Segment, error)
	// SuppressionList returns a paged result representing the client’s suppression list.
	SuppressionList(clientId string, pageSize, page int, orderBy order.SuppressionListField, direction order.Direction) (*SuppressionList, error)
	// Suppress adds the email addresses provided to the client’s suppression list.
	Suppress(clientId string, emails ...string) error
	// UnSuppress removes the email address from a client’s suppression list.
	UnSuppress(clientId string, email string) error
	Templates(clientId string) ([]*Template, error)
}
