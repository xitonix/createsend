package createsend

import (
	"fmt"
	"net/url"

	"github.com/xitonix/createsend/campaigns"
	"github.com/xitonix/createsend/internal"
)

type campaignsAPI struct {
	client internal.Client
}

func newCampaignsAPI(client internal.Client) *campaignsAPI {
	return &campaignsAPI{client: client}
}

func (a *campaignsAPI) CreateDraft(clientID string, campaign campaigns.Draft) (string, error) {
	path := fmt.Sprintf("campaigns/%s.json", url.QueryEscape(clientID))
	var result string
	err := a.client.Post(path, &result, campaign)
	if err != nil {
		return "", err
	}
	return result, nil
}
