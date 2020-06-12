package accounts

// API represents
type API interface {
	// Client returns a list of all the clients belong to the account.
	Clients() ([]*Client, error)
	// Billing returns the billing details of the account.
	Billing() (*Billing, error)
}
