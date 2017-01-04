// Package hydrasdk is a lightweight sdk for https://www.ory.am/products/hydra
// The sdk that they provide is more complete but also huge.
package hydrasdk

import (
	"net/http"
	"net/url"

	"github.com/juju/errors"
)

// Client is an oauth2 client saved on hydra database
type Client struct {
	ID                string   `json:"id" gorethink:"id"`
	Name              string   `json:"client_name" gorethink:"client_name"`
	Secret            string   `json:"client_secret,omitempty" gorethink:"client_secret"`
	RedirectURIs      []string `json:"redirect_uris" gorethink:"redirect_uris"`
	GrantTypes        []string `json:"grant_types" gorethink:"grant_types"`
	ResponseTypes     []string `json:"response_types" gorethink:"response_types"`
	Scope             string   `json:"scope" gorethink:"scope"`
	Owner             string   `json:"owner" gorethink:"owner"`
	PolicyURI         string   `json:"policy_uri" gorethink:"policy_uri"`
	TermsOfServiceURI string   `json:"tos_uri" gorethink:"tos_uri"`
	ClientURI         string   `json:"client_uri" gorethink:"client_uri"`
	LogoURI           string   `json:"logo_uri" gorethink:"logo_uri"`
	Contacts          []string `json:"contacts" gorethink:"contacts"`
	Public            bool     `json:"public" gorethink:"public"`
}

// ClientGetter is an abstraction that allows you to retrieve a specific client by their ID
type ClientGetter interface {
	Get(id string) (*Client, error)
}

// ClientManager uses hydra rest apis to retrieve clients
type ClientManager struct {
	Endpoint *url.URL
	Client   *http.Client
}

// NewClientManager returns a ClientManager connected to the hydra cluster
// it can fail if the cluster is not a valid url, or if the id and secret don't work
func NewClientManager(id, secret, cluster string) (*ClientManager, error) {
	endpoint, client, err := authenticate(id, secret, cluster)
	if err != nil {
		return nil, errors.Annotate(err, "Instantiate ClientManager")
	}

	manager := ClientManager{
		Endpoint: joinURL(endpoint, "clients"),
		Client:   client,
	}
	return &manager, nil
}

// Get queries the hydra api to retrieve a specific client by their ID.
func (m ClientManager) Get(id string) (*Client, error) {
	url := joinURL(m.Endpoint, id).String()

	var client *Client
	bind(m.Client, url, client)

	return client, nil
}
