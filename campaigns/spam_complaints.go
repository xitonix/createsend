package campaigns

import (
	"github.com/xitonix/createsend/order"
	"time"
)

// SpamComplaint represents the specific recipient and time of the spam complaint
type SpamComplaint struct {
	// Recipient represents the details of the user that made the complaint
	Recipient Recipient
	// Date the time the complaint was made
	Date time.Time
}

// SpamComplaints represents all subscribers who marked a campaign as spam
type SpamComplaints struct {
	// SpamComplaint represents the specific recipient spam complaint
	Results []SpamComplaint
	// OrderedBy the field by which the result set was ordered (email/list/date)
	OrderedBy order.Field
	// Page paginated result details
	Page order.Page
}
