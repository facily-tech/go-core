package telemetry

import (
	"context"
	"net/http"

	ddchi "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-chi/chi.v5"
	ddhttp "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const (
	// DataDogConfigPrefix is the prefix of datadog Environment.
	DataDogConfigPrefix = "DD_"

	// DataDogName is a string of type Name with holds the telemetry tool, which is "datadog".
	DataDogName Name = "datadog"
)

// Verify interface compliance.
var _ Tracer = (*DataDog)(nil)

// DataDog implements Tracer.
type DataDog struct{}

// DataDogConfig is the struct of config given to NewDataDog.
type DataDogConfig struct {
	Env     string `env:"ENV,required"`
	Service string `env:"SERVICE,required"`
	Version string
}

// NewDataDog returns a new Datadog implementation.
func NewDataDog(config DataDogConfig) *DataDog {
	tracer.Start([]tracer.StartOption{
		tracer.WithEnv(config.Env),
		tracer.WithService(config.Service),
		tracer.WithServiceVersion(config.Version),
	}...)

	return &DataDog{}
}

// Middleware add into http framework datadog tracer to track each requisition.
func (d *DataDog) Middleware(next http.Handler) http.Handler {
	return ddchi.Middleware()(next)
}

// Close datadog tracer.
func (DataDog) Close() {
	tracer.Stop()
}

// Client wraps datadog tracer into http.Client.
func (DataDog) Client(parent *http.Client) *http.Client {
	return ddhttp.WrapClient(parent)
}

// Name return Logger implementation name.
func (DataDog) Name() Name {
	return DataDogName
}

// SpanFromContext return span from context.Context.
// nolint: ireturn // it will not be changed to struct to maintain compatibility.
func (DataDog) SpanFromContext(ctx context.Context) (Span, bool) {
	rawSpan, ok := tracer.SpanFromContext(ctx)

	if !ok {
		return nil, false
	}

	rawContext := rawSpan.Context()

	return &ddSpan{
		context: ddSpanContext{
			traceID: rawContext.TraceID(),
			spanID:  rawContext.SpanID(),
		},
	}, true
}

// Client receive *http.Client and return a *http.Client with datadog tracer wrapped in it.
func Client(parent *http.Client) *http.Client {
	return ddhttp.WrapClient(parent)
}
