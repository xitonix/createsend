package transactional

import (
	"time"
)

// SmartEmailDetails represents the details of a Smart Transactional email.
type SmartEmailDetails struct {
	// ID email ID.
	ID string
	// Name email name.
	Name string
	// CreatedAt the time when the email was created.
	CreatedAt time.Time
	// Status the current status of the email.
	Status SmartEmailStatus
}
