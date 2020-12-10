package campaigns

import (
	"github.com/xitonix/createsend/order"
	"time"
)

// API is an interface that wraps campaign related operations.
//
// The API contains all the functionality you need to create, delete, send, schedule and query Campaign results.
type API interface {
	// Create creates a new campaign for the specified client based on the provided campaign data
	Create(clientID string, campaign WithURLs) (string, error)
	// Delete deletes a campaign from your account
	Delete(campaignID string) error
	// Send send campaign immediately
	Send(campaignID string, confirmationEmails ...string) error
	// ScheduleSend schedules a campaign to be sent at a specified date in the future
	SendAt(campaignID string, at time.Time, confirmationEmails ...string) error
	// Unschedule Cancels the sending of the campaign and moves it back into the drafts. If the campaign is already sent
	// or in the process of sending, this operation will fail.
	Unschedule(campaignID string) error
	// SendPreview send a test preview campaign
	SendPreview(campaignID string, recipients ...string) error
	// Summary gets the reporting summary data for the specified campaign
	Summary(campaignID string) (*Summary, error)
	// EmailClientUsage lists the email clients subscribers used to open the campaign
	EmailClientUsage(campaignID string) ([]*EmailClientUsage, error)
	// ListsAndSegments returns the lists and segments any campaign was sent to
	ListsAndSegments(campaignID string) (*ListsAndSegments, error)
	// Recipients Retrieves a paged result representing all the subscribers that a given campaign was sent to
	Recipients(campaignID string, page int, pageSize int, orderField order.Field, orderDirection order.Direction) (*Recipients, error)
	// Bounces Retrieves a paged result representing all the subscribers who bounced for a given campaign, and the type
	// of bounce (Hard = Hard Bounce, Soft = Soft Bounce)
	Bounces(campaignID string, date time.Time, page int, pageSize int, orderField order.Field, orderDirection order.Direction) (*Bounces, error)
	// Opens Retrieves a paged result representing all subscribers who opened a given campaign, including the date/time
	// and IP address from which they opened the campaign. When possible, the latitude, longitude, city, region, country
	// code, and country name as geocoded from the IP address, are also returned.
	Opens(campaignID string, date time.Time, page int, pageSize int, orderField order.Field, orderDirection order.Direction) (*Opens, error)
	// Clicks Retrieves a paged result representing all subscribers who clicked a link in a given campaign.
	// The date field is optional, opens on or after the date value specified will be returned.
	Clicks(campaignID string, date time.Time, page int, pageSize int, orderField order.Field, orderDirection order.Direction) (*Clicks, error)
	// Unsubscribes Retrieves a paged result representing all subscribers who unsubscribed from the email for a
	// given campaign.
	Unsubscribes(campaignID string, date time.Time, page int, pageSize int, orderField order.Field, orderDirection order.Direction) (*Unsubscribes, error)
	// SpamComplaints Retrieves a paged result representing all subscribers who marked the given campaign as spam,
	// including the subscriber’s list ID and the date/time they marked the campaign as spam.
	SpamComplaints(campaignID string, date time.Time, page int, pageSize int, orderField order.Field, orderDirection order.Direction) (*SpamComplaints, error)
}
