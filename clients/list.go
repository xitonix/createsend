package clients

import "time"

// List represents a list.
type List struct {
	// ID list id.
	ID string `json:"ListID"`
	// Name list name.
	Name string
}

// Subscriber represents a subscriber.
type Subscriber struct {
	// State the subscription state.
	State string
	// DateAdded the date the subscriber has subscribed in the client’s timezone.
	DateAdded time.Time
}

// SubscriberList represents a subscriber list
type SubscriberList struct {
	// List the list the subscriber belongs to.
	List
	// Subscriber the subscriber details.
	Subscriber Subscriber
}
