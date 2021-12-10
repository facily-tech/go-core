package client

import (
	"net/http"
	"time"
)

const PrefixHTTP = "HTTP_"

type Config struct {
	Timeout time.Duration `env:"TIMEOUT,required"`
}

func (c Config) SetTimeout(client *http.Client) {
	client.Timeout = c.Timeout
}

// NewHTTPClient crate a new *http.Client with Config fields set. If parent is present only config fields are changed.
func NewHTTPClient(opts ...func(*http.Client)) *http.Client {
	client := &http.Client{}
	for i := range opts {
		opts[i](client)
	}

	return client
}
