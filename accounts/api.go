package accounts

import "time"

// API represents
type API interface {
	// Client returns a list of all the clients belong to the account.
	Clients() ([]*Client, error)
	// Billing returns the billing details of the account.
	Billing() (*Billing, error)
	// Countries returns a list of all the valid countries accepted as input when a country is required, typically when creating a client.
	Countries() ([]string, error)
	// Timezones returns a list of all the valid timezones accepted as input when a timezone is required, typically when creating a client.
	Timezones() ([]string, error)
	// Now returns the current date and time in the account’s timezone.
	//
	// This is useful when, for example, you are syncing your Campaign Monitor lists with an external list,
	// allowing you to accurately determine the time on our server when you carry out the synchronization.
	Now() (time.Time, error)
}
