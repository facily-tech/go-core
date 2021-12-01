package telemetry

import (
	"net/http"

	ddchi "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-chi/chi.v5"
	ddhttp "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const DataDogConfigPrefix = "DD_"

// Verify interface compliance
var _ Tracer = (*DataDog)(nil)

type DataDog struct{}

type DataDogConfig struct {
	Env     string `env:"ENV,required"`
	Service string `env:"SERVICE,default=core-commerce-erp-integration"`
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
