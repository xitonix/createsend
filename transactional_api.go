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

func (t *transactionalAPI) SmartEmails(options ...transactional.Option) ([]*transactional.SmartEmailDetails, error) {
	ops := &transactional.Options{}
	for _, op := range options {
		op(ops)
	}

	return t.smartEmailsByStatus(ops.SmartEmailStatus(), ops.ClientID())
}

func (t *transactionalAPI) smartEmailsByStatus(status transactional.SmartEmailStatus, clientID string) ([]*transactional.SmartEmailDetails, error) {
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

	var smartEmails []internal.SmartEmailDetails
	err := t.client.Get(path, &smartEmails)
	if err != nil {
		return nil, err
	}

	result := make([]*transactional.SmartEmailDetails, len(smartEmails))
	for i, raw := range smartEmails {
		smartEmail, err := raw.ToSmartEmailDetails()
		if err != nil {
			return nil, newClientError(ErrCodeDataProcessing)
		}
		result[i] = smartEmail
	}

	return result, nil
}
