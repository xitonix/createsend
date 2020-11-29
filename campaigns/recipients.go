package campaigns

import "github.com/xitonix/createsend/order"

// Recipient represents a subscriber to a campaign
type Recipient struct {
	// EmailAddress the email address
	EmailAddress string
	// ListID the list id that the recipient is in
	ListID string
}

// RecipientLocationDetails represents the location details of a recipient
type RecipientLocationDetails struct {
	// IPAddress the IP address
	IPAddress string
	// Latitude as geocoded from the IP address
	Latitude float64
	// Longitude as geocoded from the IP address
	Longitude float64
	// City as geocoded from the IP address
	City string
	// Region as geocoded from the IP address
	Region string
	// CountryCode as geocoded from the IP address
	CountryCode string
	// CountryName as geocoded from the IP address
	CountryName string
}

// Recipients represents the subscribers to a campaign
type Recipients struct {
	// Results represents the specific recipients of a campaign
	Results []Recipient
	// OrderedBy the field by which the result set was ordered (email/list/date)
	OrderedBy order.Field
	// Page paginated result details
	Page order.Page
}
