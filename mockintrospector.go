package hydrasdk

import (
	"errors"
	"strings"
)

// IntrospecterMocker returns fictitious introspections
type IntrospecterMocker struct {
}

// NewIntrospecterMocker returns a IntrospecterMocker with a default introspection inside it
func NewIntrospecterMocker() *IntrospecterMocker {
	Mocker := IntrospecterMocker{}
	return &Mocker
}

// Introspect accepts a token in this form:
//   "userid:scope1,scope2"
// and will return an appropriate introspection
func (m IntrospecterMocker) Introspect(token string, scopes ...string) (*Introspection, error) {
	introspection := new(Introspection)

	parts := strings.Split(token, ":")
	if len(parts) != 2 {
		return nil, errors.New("The token must be in the form 'userid:scope1,scope2'")
	}

	introspection.Subject = parts[0]

	parts = strings.Split(parts[1], ",")

	introspection.Scope = strings.Join(parts, " ")

	introspection.Active = true
	for _, scope := range scopes {
		if !in(parts, scope) {
			introspection.Active = false
		}
	}

	return introspection, nil
}

func in(slice []string, a string) bool {
	for _, item := range slice {
		if a == item {
			return true
		}
	}

	return false
}
