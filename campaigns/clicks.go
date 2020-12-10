package campaigns

import (
	"github.com/xitonix/createsend/order"
	"time"
)

// ClickDetails represent specific details about each email link clicked
type ClickDetails struct {
	// Recipient represents the details of a recipient that clicked a link in the email
	Recipient
	// RecipientLocationDetails represents location details of a recipient
	RecipientLocationDetails
	// Date represents the date the open occurred
	Date time.Time
	// URL represents the URL that was clicked
	URL string
}

type Clicks struct {
	// Results represent specific details about each click
	Results []ClickDetails
	// OrderedBy the field by which the result set was ordered (email/list/date).
	OrderedBy order.Field
	// Page paginated result details
	Page order.Page
}
