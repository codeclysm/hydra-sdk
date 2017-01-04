package hydrasdk

// ClientMocker returns fictitious clients
type ClientMocker struct {
}

// NewClientMocker returns a ClientMocker with a default client inside it
func NewClientMocker() *ClientMocker {
	Mocker := ClientMocker{}
	return &Mocker
}

// Get returns the client it found inside it
func (m ClientMocker) Get(id string) (*Client, error) {
	client := new(Client)

	return client, nil
}
