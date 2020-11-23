package campaigns

import "github.com/xitonix/createsend/order"

// Recipient represents a subscriber to a campaign
type Recipient struct {
	// EmailAddress the email address
	EmailAddress string
	// ListID the list id that the recipient is in
	ListID string
}

// Recipients represents the subscribers to a campaign
type Recipients struct {
	// Results represents the specific recipients of a campaign
	Results []Recipient
	// OrderedBy the field by which the result set was ordered (email/list/date)
	OrderedBy order.Field
	// Page paginated result details
	Page order.Page
}
