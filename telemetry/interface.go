package telemetry

import "net/http"

// A Tracer has methods to help tracer instrumentation of our services.
type Tracer interface {
	// Middleware must return a new handler with cross application tracing (CAT) or distributed tracing.
	Middleware(next http.Handler) http.Handler
	// Client wraps parent with tracing capabilities, parent is modified during this process.
	Client(parent *http.Client) *http.Client

	// Close should be called when the application end.
	Close()
}
