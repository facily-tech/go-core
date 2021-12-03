package telemetry

import (
	"context"
	"net/http"

	ddchi "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-chi/chi.v5"
	ddhttp "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const DataDogConfigPrefix = "DD_"
const DataDogName Name = "datadog"

// Verify interface compliance
var _ Tracer = (*DataDog)(nil)

type DataDog struct{}

type DataDogConfig struct {
	Env     string `env:"ENV,required"`
	Service string `env:"SERVICE,required"`
	Version string
}

func NewDataDog(config DataDogConfig) *DataDog {
	tracer.Start([]tracer.StartOption{
		tracer.WithEnv(config.Env),
		tracer.WithService(config.Service),
		tracer.WithServiceVersion(config.Version),
	}...)

	return &DataDog{}
}

func (d *DataDog) Middleware(next http.Handler) http.Handler {
	return ddchi.Middleware()(next)
}

func (DataDog) Close() {
	tracer.Stop()
}

func (DataDog) Client(parent *http.Client) *http.Client {
	return ddhttp.WrapClient(parent)
}

func (DataDog) Name() Name {
	return DataDogName
}

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
