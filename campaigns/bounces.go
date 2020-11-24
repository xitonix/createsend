package campaigns

import (
	"github.com/xitonix/createsend/order"
	"time"
)

type Bounce struct {
	Recipient
	Date       time.Time
	BounceType string
	Reason     string
}

type Bounces struct {
	Results   []Bounce
	OrderedBy order.Field
	Page      order.Page
}
