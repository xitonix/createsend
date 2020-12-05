package createsend

import (
	"fmt"
	"net/url"

	"github.com/xitonix/createsend/internal"
	"github.com/xitonix/createsend/transactional"
)

type transactionalAPI struct {
	client internal.Client
}

func newTransactionalAPI(client internal.Client) *transactionalAPI {
	return &transactionalAPI{client: client}
}

func (t *transactionalAPI) SmartEmails(options ...transactional.Option) ([]*transactional.SmartEmailBasicDetails, error) {
	ops := &transactional.Options{}
	for _, op := range options {
		op(ops)
	}

	return t.smartEmailsByStatus(ops.SmartEmailStatus(), ops.ClientID())
}

func (t *transactionalAPI) SmartEmail(smartEmailID string) (*transactional.SmartEmailDetails, error) {
	path := fmt.Sprintf("transactional/smartEmail/%s", url.QueryEscape(smartEmailID))
	var smartEmail internal.SmartEmailDetails
	err := t.client.Get(path, &smartEmail)
	if err != nil {
		return nil, err
	}

	result, err := smartEmail.ToSmartEmailDetails()
	if err != nil {
		return nil, newClientError(ErrCodeDataProcessing)
	}

	return result, nil
}

func (t *transactionalAPI) smartEmailsByStatus(status transactional.SmartEmailStatus, clientID string) ([]*transactional.SmartEmailBasicDetails, error) {
	var statusParam string
	switch status {
	case transactional.UnknownSmartEmail:
		statusParam = "all"
	default:
		statusParam = status.String()
	}
	path := fmt.Sprintf("transactional/smartEmail?status=%s", statusParam)
	if clientID != "" {
		path += "&clientID=" + url.QueryEscape(clientID)
	}

	var smartEmails []internal.SmartEmailBasicDetails
	err := t.client.Get(path, &smartEmails)
	if err != nil {
		return nil, err
	}

	result := make([]*transactional.SmartEmailBasicDetails, len(smartEmails))
	for i, raw := range smartEmails {
		smartEmail, err := raw.ToSmartEmailBasicDetails()
		if err != nil {
			return nil, newClientError(ErrCodeDataProcessing)
		}
		result[i] = smartEmail
	}

	return result, nil
}
