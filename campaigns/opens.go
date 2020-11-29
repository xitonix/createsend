package campaigns

import (
	"github.com/xitonix/createsend/order"
	"time"
)

// OpenDetails represent specific details about each open
type OpenDetails struct {
	// Recipient represents the details of a recipient that opened the email
	Recipient
	// RecipientLocationDetails represents location details of a recipient
	RecipientLocationDetails
	// Date represents the date the the open occurred
	Date time.Time
}

// Opens represents all subscribers who opened a given campaign
type Opens struct {
	// Results represent specific details about each open
	Results []OpenDetails
	// OrderedBy the field by which the result set was ordered (email/list/date).
	OrderedBy order.Field
	// Page paginated result details
	Page order.Page
}
