package campaigns

import (
	"github.com/xitonix/createsend/order"
	"time"
)

// Bounce represent basic details of a bounce
type Bounce struct {
	// Recipient details of the email that bounced
	Recipient
	// Date that the bounce occurred
	Date time.Time
	// BounceType the type of bounce
	BounceType string
	// Reason more detailed information regarding the bounce
	Reason string
}

type Bounces struct {
	// Results represent basic details of a bounce
	Results []Bounce
	// OrderedBy the field by which the result set was ordered (email/list/date).
	OrderedBy order.Field
	// Page paginated result details
	Page order.Page
}
