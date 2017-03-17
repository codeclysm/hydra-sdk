package hydrasdk

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/codeclysm/introspector"
	"github.com/juju/errors"
)

// Introspector uses hydra rest apis to retrieve clients
type Introspector struct {
	AllowedEndpoint    *url.URL
	IntrospectEndpoint *url.URL
	Client             *http.Client
}

// NewIntrospector returns a Introspector connected to the hydra cluster
// it can fail if the cluster is not a valid url, or if the id and secret don't work
func NewIntrospector(id, secret, cluster string) (*Introspector, error) {
	endpoint, client, err := authenticate(id, secret, cluster)
	if err != nil {
		return nil, errors.Annotate(err, "Instantiate Introspector")
	}

	manager := Introspector{
		AllowedEndpoint:    joinURL(endpoint, "warden", "token", "allowed"),
		IntrospectEndpoint: joinURL(endpoint, "oauth2", "introspect"),
		Client:             client,
	}
	return &manager, nil
}

type req struct {
	Scopes   []string          `json:"scopes"`
	Token    string            `json:"token"`
	Resource string            `json:"resource"`
	Action   string            `json:"action"`
	Context  map[string]string `json:"context"`
}

type res struct {
	introspector.Introspection
	Allowed   bool      `json:"allowed"`
	IssuedAt  time.Time `json:"iat"`
	ExpiresAt time.Time `json:"exp"`
	Scopes    []string  `json:"scopes"`
}

// Introspect queries the endpoint with an http request. It expects that the endpoint
// implements https://tools.ietf.org/html/rfc7662
func (m *Introspector) Introspect(token string, scopes ...string) (*introspector.Introspection, error) {
	data := url.Values{
		"token": []string{token},
		"scope": []string{strings.Join(scopes, " ")},
	}

	url := m.IntrospectEndpoint.String()
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, errors.Annotatef(err, "new request for %s", url)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	var i res
	err = bind(m.Client, req, &i)
	if err != nil {
		return nil, err
	}
	i.Introspection.Scope = strings.Join(i.Scopes, " ")
	i.Introspection.IssuedAt = i.IssuedAt.Unix()
	i.Introspection.ExpiresAt = i.ExpiresAt.Unix()
	return &i.Introspection, nil

}

// Allowed calls the hydra endpoint to retrieve the info of a token and see if it has the permission to perform an action
func (m *Introspector) Allowed(token string, perm introspector.Permission, scopes ...string) (*introspector.Introspection, bool, error) {
	payload := req{
		Token:    token,
		Scopes:   scopes,
		Resource: perm.Resource,
		Action:   perm.Action,
		Context:  perm.Context,
	}

	data, err := json.Marshal(&payload)
	if err != nil {
		return nil, false, errors.Annotatef(err, "marshal payload %+v", payload)
	}

	url := m.AllowedEndpoint.String()
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, false, errors.Annotatef(err, "new request for %s", url)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(data)))

	var i res
	err = bind(m.Client, req, &i)
	if err != nil {
		return nil, false, err
	}
	i.Introspection.Scope = strings.Join(i.Scopes, " ")
	i.Introspection.IssuedAt = i.IssuedAt.Unix()
	i.Introspection.ExpiresAt = i.ExpiresAt.Unix()

	return &i.Introspection, i.Allowed, nil
}
