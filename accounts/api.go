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
	// AddAdministrator adds a new administrator to the account.
	//
	// An invitation will be sent to the new administrator via email.
	AddAdministrator(administrator Administrator) error
	// UpdateAdministrator updates the email address and/or name of an administrator.
	//
	// The first parameter is the email address of the admin whose details will be updated.
	// This is regarded as the 'old' email address.
	UpdateAdministrator(currentEmailAddress string, administrator Administrator) error
	// Administrators returns a list of all (active or invited) administrators associated with the account.
	Administrators() ([]*AdministratorDetails, error)
	// Administrator returns the details of a single administrator associated with the account.
	//
	// The parameter is the email address of the administrator whose information should be retrieved.
	Administrator(emailAddress string) (*AdministratorDetails, error)
}
