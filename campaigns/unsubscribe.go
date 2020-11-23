package campaigns

import (
	"github.com/xitonix/createsend/order"
	"time"
)

type Unsubscribe struct {
	Recipient
	Date      time.Time
	IPAddress string
}

type Unsubscribes struct {
	Results   []Unsubscribe
	OrderedBy order.Field
	Page      order.Page
}
