package log

import (
	"context"

	"github.com/facily-tech/go-core/telemetry"
)

func ExampleNewLoggerZap() {
	log, err := NewLoggerZap(ZapConfig{
		Version:           "v0.1.0",
		DisableStackTrace: false,
	})
	if err != nil {
		// panic should not be used outside func main
		// but it is not possible to return on it tests
		panic(err)
	}

	ctx := context.Background()
	log.Info(ctx, "new log created successfully")
}

func ExampleNewLoggerZap_withTracing() {
	version := "v0.1.0"

	// If is there already a tracer instance of telemetry you should use it
	// this imsplementation exists only for example how to add a log tracer on it
	tracer, err := telemetry.NewDataDog(telemetry.DataDogConfig{
		Env:     "Local",   // environment one of: local, development, homolog, production.
		Service: "Service", // application service name, it should be the same of it repository ( Github, Bitbucket, etc ).
		Version: version,   // Code Version, if there is no version control, then use Git Commit Hash instead.
	})
	if err != nil {
		// panic should not be used outside func main
		// but it is not possible to return on it tests
		panic(err)
	}

	log, err := NewLoggerZap(ZapConfig{
		Version:           version,
		DisableStackTrace: false,
		Tracer:            tracer,
	})
	if err != nil {
		// panic should not be used outside func main
		// but it is not possible to return on it tests
		panic(err)
	}

	ctx := context.Background()
	log.Info(ctx, "new log created successfully")
}
