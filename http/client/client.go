/*
Package client make easy to use an http client
*/
package client

import (
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/facily-tech/go-core/log"
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

func WithLogger(c *http.Client, log log.Logger) func(*http.Client) {
	return func(c *http.Client) {
		c.Transport = NewLogTripper(c.Transport, log)
	}
}

type HeaderTripper struct {
	http.RoundTripper
	log log.Logger
}

func NewLogTripper(parent http.RoundTripper, log log.Logger) *HeaderTripper {
	if parent == nil {
		parent = http.DefaultTransport
	}
	return &HeaderTripper{
		RoundTripper: parent,
		log:          log,
	}
}

func (rt *HeaderTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	dr, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		return nil, err
	}

	rt.log.Info(r.Context(), "http request", log.Any("dump", string(dr)))

	resp, err := rt.RoundTripper.RoundTrip(r)
	if err != nil {
		return resp, err
	}

	dresp, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		rt.log.Info(r.Context(), "http response", log.Any("dump", string(dresp)))
	}

	return resp, err
}
