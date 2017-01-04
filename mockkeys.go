package hydrasdk

import "crypto/rsa"

// KeyMocker returns fictitious Keys
type KeyMocker struct {
}

// NewKeyMocker returns a KeyMocker
func NewKeyMocker() *KeyMocker {
	Mocker := KeyMocker{}
	return &Mocker
}

// GetRSAPublic returns a fake key
func (m KeyMocker) GetRSAPublic(set string) (*rsa.PublicKey, error) {
	key := new(rsa.PublicKey)
	return key, nil
}

// GetRSAPrivate returns a previously saved key
func (m KeyMocker) GetRSAPrivate(set string) (*rsa.PrivateKey, error) {
	key := new(rsa.PrivateKey)
	return key, nil
}
