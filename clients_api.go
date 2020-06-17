package createsend

import (
	"strings"

	"github.com/xitonix/createsend/clients"
	"github.com/xitonix/createsend/internal"
)

const (
	clientsPath = "clients.json"
)

type clientsAPI struct {
	client internal.Client
}

func newClientsAPI(client internal.Client) *clientsAPI {
	return &clientsAPI{client: client}
}

func (a *clientsAPI) Create(client clients.Client) (string, error) {
	var clientId string
	err := a.client.Post(clientsPath, &clientId, client)
	return strings.Trim(clientId, `"`), err
}
