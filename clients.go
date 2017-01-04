// Package hydrasdk is a lightweight sdk for https://www.ory.am/products/hydra
// The sdk that they provide is more complete but also huge.
package hydrasdk

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/juju/errors"
	"github.com/ory-am/hydra/pkg"
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
	uri, err := url.Parse(cluster)
	if err != nil {
		return nil, errors.Annotatef(err, "parse url %s", cluster)
	}

	credentials := clientcredentials.Config{
		ClientID:     id,
		ClientSecret: secret,
		TokenURL:     pkg.JoinURL(uri, "oauth2/token").String(),
		Scopes:       []string{"hydra"},
	}

	ctx := context.Background()
	_, err = credentials.Token(ctx)
	if err != nil {
		return nil, errors.Annotatef(err, "connect to cluster %s", cluster)
	}

	manager := ClientManager{
		Endpoint: joinURL(uri, "clients"),
		Client:   credentials.Client(ctx),
	}
	return &manager, nil
}

// Get queries the hydra api to retrieve a specific client by their ID.
func (m ClientManager) Get(id string) (*Client, error) {
	url := pkg.JoinURL(m.Endpoint, id).String()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Annotatef(err, "new request for %s", url)
	}

	resp, err := m.Client.Do(req)
	if err != nil {
		return nil, errors.Annotatef(err, "execute request %+v", req)
	}
	defer resp.Body.Close()

	var client *Client

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.Errorf("Expected status code %d, got %d.\n%s\n", http.StatusOK, resp.StatusCode, body)
	} else if err := json.NewDecoder(resp.Body).Decode(client); err != nil {
		return nil, errors.Annotatef(err, "decode json %s", resp.Body)
	}

	return client, nil
}

func joinURL(u *url.URL, args ...string) (ep *url.URL) {
	ep = copyURL(u)
	ep.Path = path.Join(append([]string{ep.Path}, args...)...)
	return ep
}

func copyURL(u *url.URL) *url.URL {
	a := new(url.URL)
	*a = *u
	return a
}
