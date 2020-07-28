package clients

// ContactDetails represents the contact details of a Person within a Client.
type ContactDetails struct {
	// Name contact's name.
	Name string
	// EmailAddress contact's email address.
	EmailAddress string
	// AccessLevel access level
	AccessLevel int
	// Username username.
	Username string
}

// ClientDetails represents Client details.
type ClientDetails struct {
	// APIKey Client API key
	APIKey string
	// ID client id
	ID string
	// Company company name
	Company string
	// Country country
	Country string
	// Timezone timezone
	Timezone string
	// Contact the contact details and access level of the Person within the client.
	//
	// If there are multiple Persons in this Client, or no Persons at all, this will be nil.
	Contact *ContactDetails
	// Billing is the Client's billing details.
	//
	// If the authenticated user is NOT authorised to see the billing details, the value will be nil.
	Billing *BillingDetails
}
