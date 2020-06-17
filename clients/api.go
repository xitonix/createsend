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
	//
	// If there is only one Person in this Client, their Email Address, Contact Name and Access Details are returned.
	// If there are multiple Persons in this Client, or no Persons at all, these fields are omitted in the response.
	Get(clientId string) (*ClientDetails, error)
}
