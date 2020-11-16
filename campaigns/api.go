package campaigns

import (
	"github.com/xitonix/createsend/campaigns/orderfield"
	"github.com/xitonix/createsend/order"
	"time"
)

type CampaignForCreation struct {
	Name       string
	Subject    string
	FromName   string
	FromEmail  string
	ReplyTo    string
	ListIds    []string
	SegmentIds []string
}

type CampaignFromUrl struct {
	CampaignForCreation
	HtmlUrl string
	// TextUrl is optional and if provided as either null or an
	// empty string, text content for the campaign will be generated from
	// the HTML content.
	TextUrl string
}

type EditableField struct {
	Content string
	Alt     string
	Href    string
}

type Repeater struct {
	Items []struct {
		Layout      []string
		Singlelines []EditableField
		Multilines  []EditableField
		Images      []EditableField
	}
}

type CampaignFromTemplate struct {
	CampaignForCreation
	TemplateID      string
	TemplateContent struct {
		Singlelines []EditableField
		Multilines  []EditableField
		Images      []EditableField
		Repeaters   []Repeater
	}
}

type CampaignSummary struct {
	Recipients        int
	TotalOpened       int
	Clicks            int
	Unsubscribed      int
	Bounced           int
	UniqueOpened      int
	SpamComplaints    int
	WebVersionURL     string
	WebVersionTextURL string
	WorldviewURL      string
	Forwards          int
	Likes             int
	Mentions          int
}

type EmailClientUsage struct {
	Client      string
	Version     string
	Percentage  float32
	Subscribers int
}

type List struct {
	ListID string
	Name   string
}

type Segment struct {
	ListID    string
	SegmentID string
	Title     string
}

type ListsAndSegments struct {
	Lists    []List
	Segments []Segment
}

type Result struct {
	EmailAddress string
	ListID       string
}

type ListRequestDetails struct {
	ResultsOrderedBy     orderfield.OrderField
	OrderDirection       order.Direction
	PageNumber           int
	PageSize             int
	RecordsOnThisPage    int
	TotalNumberOfRecords int
	NumberOfPages        int
}

type Recipients struct {
	Results []Result
	ListRequestDetails
}

type Bounces struct {
	Results []struct {
		Result
		BounceType string
		Date       time.Time
		Reason     string
	}
	ListRequestDetails
}

type CampaignRecipientActions struct {
	Results []struct {
		Result
		Date        time.Time
		IPAddress   string
		Latitude    float64
		Longitude   float64
		City        string
		Region      string
		CountryCode string
		CountryName string
	}
	ListRequestDetails
}

type Unsubscribes struct {
	Results []struct {
		Result
		EmailAddress string
		ListID       string
		Date         time.Time
		IPAddress    string
	}
	ListRequestDetails
}

type SpamComplaints struct {
	Results []struct {
		Result
		Date time.Time
	}
	ListRequestDetails
}

// API is an interface that wraps campaign related operations.
//
// The API contains all the functionality you need to create, delete, send, schedule and query Campaign results.
type API interface {
	// Create creates a new campaign for the specified client based on the provided campaign data.
	Create(clientID string, campaign CampaignFromUrl) (string, error)
	// Delete deletes a campaign from your account. For draft and scheduled campaigns (prior to the time of scheduling),
	// this will prevent the campaign from sending. If the campaign is already sent or in the process of sending, it
	// will remove the campaign from the account.
	Delete(campaignID string) error
	// SendImmediately send campaign immediately
	//
	// confirmationEmails is an email address (or a maximum of five comma-separated email addresses) to which a
	// confirmation email will be sent once your campaign has been sent.
	SendImmediately(campaignID string, confirmationEmails string) error
	// ScheduleSend schedules a campaign to be sent at a specified date in the future
	//
	// confirmationEmails is an email address (or a maximum of five comma-separated email addresses) to which a
	// confirmation email will be sent once your campaign has been sent.
	//
	// date is the future date the campaign should be scheduled to be sent. This date should be in the client's timezone.
	ScheduleSend(campaignID string, confirmationEmails string, date time.Time) error
	// Unschedule Cancels the sending of the campaign and moves it back into the drafts. If the campaign is already sent
	// or in the process of sending, this operation will fail.
	Unschedule(campaignID string) error
	// Test send a test preview campaign
	//
	// PreviewRecipients the recipients to send the preview to
	Test(campaignID string, previewRecipients []string) error
	// Summary gets the reporting summary data for the specified campaign
	Summary(campaignID string) (CampaignSummary, error)
	// EmailClientUsage lists the email clients subscribers used to open the campaign
	EmailClientUsage(campaignID string) ([]EmailClientUsage, error)
	// ListsAndSegments returns the lists and segments any campaign was sent to
	ListsAndSegments(campaignID string) (ListsAndSegments, error)
	// Recipients Retrieves a paged result representing all the subscribers that a given  campaign was sent to. This
	// includes their email address and the ID of the list they are a member of.
	// You have complete control over how results should be returned including page size, sort order and sort direction
	//
	// campaignId The id of the campaign to retrieve results for
	//
	// page The results page to retrieve. Default: 1.
	//
	// pageSize The number of records to retrieve per results page
	//
	// orderField The field which should be used to order the results
	//
	// orderDirection The direction in which results should be ordered
	Recipients(campaignID string, page int, pageSize int, orderField orderfield.OrderField, orderDirection order.Direction) (Recipients, error)
	// Bounces Retrieves a paged result representing all the subscribers who bounced for a given campaign, and the type
	// of bounce (Hard = Hard Bounce, Soft = Soft Bounce).
	// The date field is optional, opens on or after the date value specified will be returned.
	Bounces(campaignID string, date time.Time, page int, pageSize int, orderField orderfield.OrderField, orderDirection order.Direction) (Bounces, error)
	// Opens Retrieves a paged result representing all subscribers who opened a given campaign, including the date/time
	// and IP address from which they opened the campaign. When possible, the latitude, longitude, city, region, country
	// code, and country name as geocoded from the IP address, are also returned. You have complete control over how
	// results should be returned including page size, sort order and sort direction. The date field is optional, opens
	// on or after the date value specified will be returned.
	Opens(campaignID string, date time.Time, page int, pageSize int, orderField orderfield.OrderField, orderDirection order.Direction) (CampaignRecipientActions, error)
	// Clicks Retrieves a paged result representing all subscribers who clicked a link in a given campaign.
	// The date field is optional, opens on or after the date value specified will be returned.
	Clicks(campaignID string, date time.Time, page int, pageSize int, orderField orderfield.OrderField, orderDirection order.Direction) (CampaignRecipientActions, error)
	// Unsubscribes Retrieves a paged result representing all subscribers who unsubscribed from the email for a
	// given campaign.
	// The date field is optional, opens on or after the date value specified will be returned.
	Unsubscribes(campaignID string, date time.Time, page int, pageSize int, orderField orderfield.OrderField, orderDirection order.Direction) (Unsubscribes, error)
	// SpamComplaints Retrieves a paged result representing all subscribers who marked the given campaign as spam,
	// including the subscriber’s list ID and the date/time they marked the campaign as spam.
	// You have complete control over how results should be returned including page size, sort order and sort direction.
	// The date field is optional, opens on or after the date value specified will be returned.
	SpamComplaints(campaignID string, date time.Time, page int, pageSize int, orderField orderfield.OrderField, orderDirection order.Direction) (SpamComplaints, error)
}
