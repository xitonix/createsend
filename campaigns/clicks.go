package campaigns

import (
	"github.com/xitonix/createsend/order"
	"time"
)

type ClickDetails struct {
	Recipient
	RecipientLocationDetails
	Date time.Time
	URL  string
}

type Clicks struct {
	Results   []ClickDetails
	OrderedBy order.Field
	Page      order.Page
}
