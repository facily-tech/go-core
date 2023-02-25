/*
Package client make easy to use an http client
*/
package client

import (
	"net/http"
	"net/http/httputil"
	"regexp"
	"strconv"
	"time"

	utils "github.com/facily-tech/go-core/http/utils"
	"github.com/facily-tech/go-core/log"
	"github.com/facily-tech/go-core/telemetry"
	"github.com/pkg/errors"
)

// PrefixHTTP is the env prefix of Config environment variables.
const PrefixHTTP = "HTTP_"

// Config of the http client.
type Config struct {
	Timeout                time.Duration `env:"TIMEOUT,required"`
	RoundtripperStatusCode []string      `env:"ROUNDTRIPPER_STATUSCODE,default=2.."`
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

// WithLogger wrap http client transport with log.
// If response code does not match any acceptedStatusCode, response body is logged.
func WithLogger(c *http.Client, log log.Logger, acceptedStatusCode []string) {
	c.Transport = NewLogTripper(c.Transport, log, acceptedStatusCode)
}

// HeaderTripper wrap http.RoundTripper to enrich with log. This way all outgoing
// requests (even bodies) will be logged.
type HeaderTripper struct {
	http.RoundTripper
	log log.Logger

	responseStatusAcceptList []string
}

// NewLogTripper creates a new HeaderTripper using parent http.RoundTripper with log logger and will log responses
// bodies if no acceptedSatusCodes match (regex are accepted).
func NewLogTripper(parent http.RoundTripper, log log.Logger, acceptedStatusCodes []string) *HeaderTripper {
	if parent == nil {
		parent = http.DefaultTransport
	}

	return &HeaderTripper{
		RoundTripper: parent,
		log:          log,

		responseStatusAcceptList: acceptedStatusCodes,
	}
}

// RoundTrip refer to http.RoundTripper.
func (rt *HeaderTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	dr, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rt.log.Info(r.Context(), "http request", log.Any("dump", string(dr)))

	resp, err := rt.RoundTripper.RoundTrip(r)
	if err != nil {
		return resp, errors.WithStack(err)
	}

	dresp, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return resp, errors.WithStack(err)
	}
	if !rt.accept(resp.StatusCode) {
		utils.StatusLevel(rt.log, resp.StatusCode, utils.ClientMode)(r.Context(), "http response", log.Any("dump", string(dresp)))
	}

	return resp, errors.WithStack(err)
}

func (rt *HeaderTripper) accept(statusCode int) bool {
	code := strconv.Itoa(statusCode)
	matches := 0
	for _, v := range rt.responseStatusAcceptList {
		validStatus := regexp.MustCompile(v)
		if validStatus.MatchString(code) {
			matches++
		}
	}

	return matches > 0
}
