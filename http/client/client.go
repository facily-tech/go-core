/*
Package client make easy to use an http client
*/
package client

import (
	"net/http"
	"time"

	"github.com/facily-tech/go-core/telemetry"
)

// PrefixHTTP is the env prefix of Config environment variables.
const PrefixHTTP = "HTTP_"

// Config of the http client.
type Config struct {
	Timeout time.Duration `env:"TIMEOUT,required"`
}

// SetTimeout set request timeout.
func (c Config) SetTimeout(client *http.Client) {
	client.Timeout = c.Timeout
}

// NewHTTPClient crate a new *http.Client with Config fields set. If parent is present only config fields are changed.
func NewHTTPClient(tracer telemetry.Tracer, opts ...func(*http.Client)) *http.Client {
	client := &http.Client{}
	for i := range opts {
		opts[i](client)
	}

	return tracer.Client(client)
}
