package coincross

import (
	"net/http"
)

type newClientFunc func(apikey, secret string, transport *http.Transport) Client

var registry = make(map[string]newClientFunc)

// Register a new client to the default registry.
func Register(name string, newfunc newClientFunc) {
	registry[name] = newfunc
}

// New creates a client instance with given parameters.
func New(name string, apikey, secret string, transport *http.Transport) Client {
	newfunc, ok := registry[name]
	if ok {
		return newfunc(apikey, secret, transport)
	} else {
		return nil
	}
}

// List returns all registered client types.
func List() []string {
	exchanges := make([]string, 0)
	for name := range registry {
		exchanges = append(exchanges, name)
	}
	return exchanges
}
