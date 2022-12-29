/*
Package client make easy to use an http client
*/
package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/facily-tech/go-core/env"
	"github.com/facily-tech/go-core/log"
	"github.com/facily-tech/go-core/telemetry"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type TelemetryMock struct{}

// Middleware must return a new handler with cross application tracing (CAT) or distributed tracing.
func (m *TelemetryMock) Middleware(next http.Handler) http.Handler {
	panic("not implemented") // TODO: Implement
}

// Client wraps parent with tracing capabilities, parent is modified during this process.
func (m *TelemetryMock) Client(parent *http.Client) *http.Client {
	return &http.Client{}
}

// Close should be called when the application end.
func (m *TelemetryMock) Close() {
	panic("not implemented") // TODO: Implement
}

// Return the Name of Which implementation is using ex: DataDog, NewRelic
func (m *TelemetryMock) Name() telemetry.Name {
	panic("not implemented") // TODO: Implement
}

// Get SpanFromContext given
func (m *TelemetryMock) SpanFromContext(ctx context.Context) (telemetry.Span, bool) {
	panic("not implemented") // TODO: Implement
}

func TestWithLogger(t *testing.T) {
	t.Run("success, request info", func(t *testing.T) {
		c := &http.Client{}
		logMock := log.NewMockLogger(gomock.NewController(t))
		logMock.EXPECT().Info(gomock.Any(), "http request", gomock.Any()).Times(1)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		defer srv.Close()

		WithLogger(c, logMock, []string{"2.."})

		resp, err := c.Get(srv.URL)
		assert.NoError(t, err)
		defer resp.Body.Close()

	})

	t.Run("success, response body", func(t *testing.T) {
		c := &http.Client{}
		logMock := log.NewMockLogger(gomock.NewController(t))
		logMock.EXPECT().Info(gomock.Any(), "http request", gomock.Any()).Times(1)
		logMock.EXPECT().Info(gomock.Any(), "http response", gomock.Any()).Times(1)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusInternalServerError) }))
		defer srv.Close()

		WithLogger(c, logMock, []string{"2.."})

		resp, err := c.Get(srv.URL)
		assert.NoError(t, err)
		defer resp.Body.Close()

	})

	t.Run("we should not call response log because we are accepting 5xx status", func(t *testing.T) {
		c := &http.Client{}
		logMock := log.NewMockLogger(gomock.NewController(t))
		logMock.EXPECT().Info(gomock.Any(), "http request", gomock.Any()).Times(1)
		logMock.EXPECT().Info(gomock.Any(), "http response", gomock.Any()).Times(0)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusInternalServerError) }))
		defer srv.Close()

		WithLogger(c, logMock, []string{"[235].."})

		resp, err := c.Get(srv.URL)
		assert.NoError(t, err)
		defer resp.Body.Close()

	})

	t.Run("complete example with trigger, roundtripper logger and accept code", func(t *testing.T) {
		httpClient := NewHTTPClient(&TelemetryMock{})
		t.Setenv("HTTP_TIMEOUT", "10s")

		config := Config{}
		err := env.LoadEnv(context.Background(), &config, PrefixHTTP)
		assert.NoError(t, err)

		config.SetTimeout(httpClient)

		logMock := log.NewMockLogger(gomock.NewController(t))
		logMock.EXPECT().Info(gomock.Any(), "http request", gomock.Any()).Times(1)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		defer srv.Close()

		WithLogger(httpClient, logMock, config.RoundtripperStatusCode)

		resp, err := httpClient.Get(srv.URL)
		assert.NoError(t, err)
		defer resp.Body.Close()

	})
}
