package campaigns

import "time"

// API is an interface that wraps Campaign related operations.
//
// The API contains all the functionality you need to manage your Campaigns.
type API interface {
	// CreateDraft creates a draft campaign ready to be tested as a preview or sent under the specified client.
	//
	// You may optionally specify the Text URL if you want to specify text content for the campaign.
	// If you don’t specify Text URL or if the Text URL is left empty, the text content for the campaign will
	// be automatically generated from the HTML content.
	//
	// If you are using the Segments, remove the Lists from your request.
	CreateDraft(clientID string, campaign Draft) (string, error)
	// Send sends a draft campaign.
	Send(draftCampaignID string, confirmationEmails ...string) error
	// SendAt sends a draft campaign at the specified time in the future.
	//
	// The send date should be in the client's timezone.
	SendAt(draftCampaignID string, at time.Time, confirmationEmails ...string) error
}
