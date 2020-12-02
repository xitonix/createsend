package internal

import (
	"github.com/araddon/dateparse"

	"github.com/xitonix/createsend/transactional"
)

// SmartEmailDetails raw smart email details model.
type SmartEmailDetails struct {
	// ID email ID.
	ID string
	// Name email name.
	Name string
	// CreatedAt the time when the email was created.
	CreatedAt string
	// Status the current status of the email.
	Status transactional.SmartEmailStatus
}

// ToSmartEmailDetails converts the raw model to a new createsend model.
func (s *SmartEmailDetails) ToSmartEmailDetails() (*transactional.SmartEmailDetails, error) {
	if s == nil {
		return nil, nil
	}

	date, err := dateparse.ParseAny(s.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &transactional.SmartEmailDetails{
		ID:        s.ID,
		Name:      s.Name,
		CreatedAt: date,
		Status:    s.Status,
	}, nil
}
