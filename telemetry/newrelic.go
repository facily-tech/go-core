package telemetry

import (
	"context"
	"net/http"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
)

// Tracer verify interface compliance
// On Hold New Relic implementation will.
var _ Tracer = (*NewRelic)(nil)

// NewRelicName is a string of type Name, which is "newrelic".
const NewRelicName Name = "newrelic"

// NewRelic implements Tracer.
type NewRelic struct {
	app *newrelic.Application
}

// NewNewRelic Create a new NewRelic instance.
func NewNewRelic() (*NewRelic, error) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigFromEnvironment(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "can't initialize new relic telemetry")
	}

	return &NewRelic{app: app}, nil
}

// Middleware add into http framework request tracing.
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

// Client wraps into http client Newrelic span.
func (NewRelic) Client(parent *http.Client) *http.Client {
	if parent.Transport == nil {
		parent.Transport = http.DefaultTransport
	}
	parent.Transport = newrelic.NewRoundTripper(parent.Transport)

	return parent
}

// Close does nothing because it is not required to close connection with NewRelic Agent.
func (relic *NewRelic) Close() {
	relic.app.Shutdown(time.Minute)
}

// Name returns the name of telemetry implementation used, this case is newrelic.
func (NewRelic) Name() Name {
	return NewRelicName
}

// SpanFromContext does nothing it is a mock method.
func (relic *NewRelic) SpanFromContext(ctx context.Context) (Span, bool) {
	tx := newrelic.FromContext(ctx)
	if tx == nil {
		return nil, false
	}

	return &txn{tx}, true
}

type txn struct {
	tx *newrelic.Transaction
}

func (t *txn) Context() SpanContext {
	return t
}

func (t *txn) SpanID() interface{} {
	return t.tx.GetLinkingMetadata().SpanID
}

func (t *txn) TraceID() interface{} {
	return t.tx.GetLinkingMetadata().TraceID
}

func (t *txn) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"traceID": t.TraceID(),
		"spanID":  t.SpanID(),
	}
}
