package transactional

// API is an interface that covers Transactional emails API.
type API interface {
	// SmartEmails returns a list of all smart transactional emails.
	//
	// Use WithClientID and WithSmartEmailStatus to filter by client ID and status respectively.
	// If you are an agency using an account API key or OAuth, you will be required to specify the client.
	// This is not required if you use a client-specific API key.
	SmartEmails(options ...Option) ([]*SmartEmailBasicDetails, error)
	// SmartEmail returns the details of a smart transactional email.
	SmartEmail(smartEmailID string) (*SmartEmailDetails, error)
}
