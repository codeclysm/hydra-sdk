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
)

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

// bind does a get request and binds the body to the given interface
func bind(client *http.Client, req *http.Request, o interface{}) error {
	resp, err := client.Do(req)
	if err != nil {
		return errors.Annotatef(err, "execute request %+v", req)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.Errorf("Expected status code %d, got %d.\n%s\n", http.StatusOK, resp.StatusCode, body)
	} else if err := json.NewDecoder(resp.Body).Decode(o); err != nil {
		return errors.Annotatef(err, "decode json %s", resp.Body)
	}
	return nil
}

func authenticate(id, secret, cluster string) (*url.URL, *http.Client, error) {
	uri, err := url.Parse(cluster)
	if err != nil {
		return nil, nil, errors.Annotatef(err, "parse url %s", cluster)
	}
	credentials := clientcredentials.Config{
		ClientID:     id,
		ClientSecret: secret,
		TokenURL:     joinURL(uri, "oauth2/token").String(),
		Scopes:       []string{"hydra"},
	}

	ctx := context.Background()
	_, err = credentials.Token(ctx)
	if err != nil {
		return nil, nil, errors.Annotatef(err, "connect to cluster %s", cluster)
	}
	return uri, credentials.Client(ctx), nil
}
