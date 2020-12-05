package internal

import (
	"github.com/araddon/dateparse"

	"github.com/xitonix/createsend/transactional"
)

// SmartEmailBasicDetails raw smart email basic details model.
type SmartEmailBasicDetails struct {
	// ID email ID.
	ID string
	// Name email name.
	Name string
	// CreatedAt the time when the email was created.
	CreatedAt string
	// Status the current status of the email.
	Status transactional.SmartEmailStatus
}

// ToSmartEmailBasicDetails converts the raw model to a new createsend model.
func (s *SmartEmailBasicDetails) ToSmartEmailBasicDetails() (*transactional.SmartEmailBasicDetails, error) {
	if s == nil {
		return nil, nil
	}

	date, err := dateparse.ParseAny(s.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &transactional.SmartEmailBasicDetails{
		ID:        s.ID,
		Name:      s.Name,
		CreatedAt: date,
		Status:    s.Status,
	}, nil
}
