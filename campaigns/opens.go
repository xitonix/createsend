package campaigns

import (
	"github.com/xitonix/createsend/order"
	"time"
)

type OpenDetails struct {
	Recipient
	RecipientLocationDetails
	Date time.Time
}

type Opens struct {
	Results   []OpenDetails
	OrderedBy order.Field
	Page      order.Page
}
