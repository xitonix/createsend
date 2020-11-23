package internal

import (
	"github.com/araddon/dateparse"

	"github.com/xitonix/createsend/clients"
	"github.com/xitonix/createsend/order"
)

// SuppressionDetails represents a suppression list item.
type SuppressionDetails struct {
	// SuppressionReason reason for suppression.
	SuppressionReason string
	// EmailAddress the suppressed email address.
	EmailAddress string
	// Date the date when the email address has been added to the suppression list.
	Date string
	// State the state of the suppressed email address.
	State string
}

// SuppressionList represents client suppression list.
type SuppressionList struct {
	// Entries the list of suppressed email addresses.
	Results []*SuppressionDetails
	// OrderedBy the field by which the result set was ordered (email/date).
	ResultsOrderedBy order.SuppressionListField
	// OrderDirection the order in which the results were sorted.
	OrderDirection order.Direction
	// PageNumber the current page number.
	PageNumber int
	// PageSize the page size.
	PageSize int
	// RecordsOnThisPage the number of records on this page.
	RecordsOnThisPage int
	// TotalNumberOfRecords the total number of records.
	TotalNumberOfRecords int
	// NumberOfPages the total number of pages.
	NumberOfPages int
}

// ToSuppressionList converts the raw model to a new createsend model.
func (s *SuppressionList) ToSuppressionList() (*clients.SuppressionList, error) {
	output := &clients.SuppressionList{
		Entries:   make([]*clients.SuppressionDetails, len(s.Results)),
		OrderedBy: s.ResultsOrderedBy,
		Page: order.Page{
			OrderDirection: s.OrderDirection,
			Number:         s.PageNumber,
			Size:           s.PageSize,
			Records:        s.RecordsOnThisPage,
			Total:          s.TotalNumberOfRecords,
			NumberOfPages:  s.NumberOfPages,
		},
	}
	for i, entry := range s.Results {
		date, err := dateparse.ParseAny(entry.Date)
		if err != nil {
			return nil, err
		}
		output.Entries[i] = &clients.SuppressionDetails{
			Reason:       entry.SuppressionReason,
			EmailAddress: entry.EmailAddress,
			Date:         date,
			State:        entry.State,
		}
	}
	return output, nil
}
