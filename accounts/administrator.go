package accounts

// Administrator represents an account administrator.
type Administrator struct {
	// EmailAddress email address.
	EmailAddress string
	// Name name.
	Name string
}

// AdministratorDetails represents an account administrator details.
type AdministratorDetails struct {
	Administrator
	// Status invitation status.
	Status string
}
