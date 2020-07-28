package internal

import (
	"github.com/araddon/dateparse"

	"github.com/xitonix/createsend/clients"
)

// SubscriberList represents a raw subscriber list.
type SubscriberList struct {
	// ListID list Id.
	ListID              string
	// ListName list name.
	ListName            string
	// SubscriberState subscriber status.
	SubscriberState     string
	// DateSubscriberAdded date the subscriber was added to the list.
	DateSubscriberAdded string
}

// ToSubscriberList converts the raw model to a new createsend model.
func (s *SubscriberList) ToSubscriberList() (*clients.SubscriberList, error) {
	if s == nil {
		return nil, nil
	}

	date, err := dateparse.ParseAny(s.DateSubscriberAdded)
	if err != nil {
		return nil, err
	}

	return &clients.SubscriberList{
		List: clients.List{
			Id:   s.ListID,
			Name: s.ListName,
		},
		Subscriber: clients.Subscriber{
			State:     s.SubscriberState,
			DateAdded: date,
		},
	}, nil
}
