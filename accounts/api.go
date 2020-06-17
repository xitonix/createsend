package accounts

import "time"

// API is an interface that wraps account related operations.
//
// The API gives you access to core account information such as the clients available in your account and helper
// procedures when creating a client including available countries, time zones, the current date and
// time in your account etc.
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
	// GetAdministrators returns a list of all (active or invited) administrators associated with the account.
	GetAdministrators() ([]*AdministratorDetails, error)
	// GetAdministrator returns the details of a single administrator associated with the account.
	//
	// The parameter is the email address of the administrator whose information should be retrieved.
	GetAdministrator(emailAddress string) (*AdministratorDetails, error)
	// DeleteAdministrator changes the status of an active administrator, defined by the email address, to deleted.
	//
	// They will no longer be able to log into their account.
	DeleteAdministrator(emailAddress string) error
	// SetAsPrimaryContact sets the primary contact of the account to be the administrator with the specified email address.
	SetAsPrimaryContact(emailAddress string) error
	// GetPrimaryContact returns the email address of the administrator who is selected as the primary contact for the account.
	GetPrimaryContact() (string, error)
	// NewEmbeddedSession initiates a new login session for the member with the specified email address and returns the session URL.
	//
	// This method will return a single use URL which will create the login session.
	// This is usually used as the source of an iframe for embedding Campaign Monitor within your own application.
	NewEmbeddedSession(session EmbeddedSession) (string, error)
}
