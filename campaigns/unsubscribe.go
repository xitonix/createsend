package campaigns

import (
	"github.com/xitonix/createsend/order"
	"time"
)

// Unsubscribe represents specific details about a recipient that unsubscribed
type Unsubscribe struct {
	Recipient
	// Date that the unsubscribe occurred
	Date time.Time
	// IPAddress where the unsubscribe originated from
	IPAddress string
}

// Unsubscribes represents recipients that have unsubscribed from a campaign
type Unsubscribes struct {
	// Results specific details about a recipient that unsubscribed
	Results []Unsubscribe
	// OrderedBy the field by which the result set was ordered (email/list/date).
	OrderedBy order.Field
	// Page paginated result details
	Page order.Page
}
