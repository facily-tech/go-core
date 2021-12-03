package telemetry

import (
	"context"
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
)

// Verify interface compliance
// On Hold New Relic implementation will
var _ Tracer = (*NewRelic)(nil)

const NewRelicName Name = "newrelic"

type NewRelic struct {
	app *newrelic.Application
}

func NewNewRelic() (*NewRelic, error) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigFromEnvironment(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "can't initialize new relic telemetry")
	}
	return &NewRelic{app: app}, nil
}

func (relic *NewRelic) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		txn := relic.app.StartTransaction(r.Method + " " + r.URL.Path)
		defer txn.End()

		w = txn.SetWebResponse(w)
		txn.SetWebRequestHTTP(r)
		r = newrelic.RequestWithTransactionContext(r, txn)
		next.ServeHTTP(w, r)
	})
}

type HTTPClient struct {
	Client *http.Client
}

func (NewRelic) Client(parent *http.Client) *http.Client {
	if parent.Transport == nil {
		parent.Transport = http.DefaultTransport
	}
	parent.Transport = newrelic.NewRoundTripper(parent.Transport)
	return parent
}

func (NewRelic) Close() {
}

func (NewRelic) Name() Name {
	return NewRelicName
}

func (NewRelic) SpanFromContext(ctx context.Context) (Span, bool) {
	return nil, false
}
