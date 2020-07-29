package clients

// PersonBasicDetails represents a person's basic details
type PersonBasicDetails struct {
	// EmailAddress the email address.
	EmailAddress string
	// Name the name.
	Name string
	// AccessLevel access level.
	AccessLevel int
}

// PersonDetails represents a person's details.
type PersonDetails struct {
	PersonBasicDetails
	// Status the person's status (eg. Active)
	Status string
}

// Person represents a person.
type Person struct {
	PersonBasicDetails
	// Password the person's password.
	Password string
}
