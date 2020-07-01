package internal

import (
	"github.com/araddon/dateparse"

	"github.com/xitonix/createsend/clients"
)

type SubscriberList struct {
	ListID              string
	ListName            string
	SubscriberState     string
	DateSubscriberAdded string
}

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
