package clients

import (
	"time"

	"github.com/xitonix/createsend/order"
)

// SuppressionDetails represents a suppression list item.
type SuppressionDetails struct {
	// Reason reason for suppression.
	Reason string
	// EmailAddress the suppressed email address.
	EmailAddress string
	// Date the date when the email address has been added to the suppression list.
	Date time.Time
	// State the state of the suppressed email address.
	State string
}

// SuppressionList represents client suppression list.
type SuppressionList struct {
	// Entries the list of suppressed email addresses.
	Entries []*SuppressionDetails
	// OrderedBy the field by which the result set was ordered (email/date).
	OrderedBy order.SuppressionListField
	// Page paginated result details
	Page order.Page
}
