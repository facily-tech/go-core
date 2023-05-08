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
	panic("not implemented")
}

// Client wraps parent with tracing capabilities, parent is modified during this process.
func (m *TelemetryMock) Client(parent *http.Client) *http.Client {
	return &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}
}

// Close should be called when the application end.
func (m *TelemetryMock) Close() {
	panic("not implemented")
}

// Return the Name of Which implementation is using ex: DataDog, NewRelic.
func (m *TelemetryMock) Name() telemetry.Name {
	panic("not implemented")
}

// Get SpanFromContext.
func (m *TelemetryMock) SpanFromContext(ctx context.Context) (telemetry.Span, bool) {
	panic("not implemented")
}

func TestWithLogger(t *testing.T) {
	t.Run("success, request info", func(t *testing.T) {
		c := &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       0,
		}

		logMock := log.NewMockLogger(gomock.NewController(t))
		logMock.EXPECT().Info(gomock.Any(), "http request", gomock.Any()).Times(1)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		defer srv.Close()

		WithLogger(c, logMock, []string{"2.."})

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
		assert.NoError(t, err)

		resp, err := c.Do(req) //nolint:bodyclose // false positive our request body is nil
		assert.NoError(t, err)
		defer closeHelper(t, resp.Body)

	})

	t.Run("success, response info body", func(t *testing.T) {
		c := &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       0,
		}
		logMock := log.NewMockLogger(gomock.NewController(t))
		logMock.EXPECT().Info(gomock.Any(), "http request", gomock.Any()).Times(1)
		logMock.EXPECT().Info(gomock.Any(), "http response", gomock.Any()).Times(1)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))
		defer srv.Close()

		WithLogger(c, logMock, []string{"1.."})

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL+"/gold", nil)
		assert.NoError(t, err)

		resp, err := c.Do(req) //nolint:bodyclose // false positive our request body is nil
		assert.NoError(t, err)
		defer closeHelper(t, resp.Body)
	})

	t.Run("success, response warning body", func(t *testing.T) {
		c := &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       0,
		}
		logMock := log.NewMockLogger(gomock.NewController(t))
		logMock.EXPECT().Info(gomock.Any(), "http request", gomock.Any()).Times(1)
		logMock.EXPECT().Warn(gomock.Any(), "http response", gomock.Any()).Times(1)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusMultipleChoices) }))
		defer srv.Close()

		WithLogger(c, logMock, []string{"2.."})

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
		assert.NoError(t, err)

		resp, err := c.Do(req) //nolint:bodyclose // false positive our request body is nil
		assert.NoError(t, err)
		defer closeHelper(t, resp.Body)
	})

	t.Run("success, response error body", func(t *testing.T) {
		c := &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       0,
		}
		logMock := log.NewMockLogger(gomock.NewController(t))
		logMock.EXPECT().Info(gomock.Any(), "http request", gomock.Any()).Times(1)
		logMock.EXPECT().Error(gomock.Any(), "http response", gomock.Any()).Times(1)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusBadRequest) }))
		defer srv.Close()

		WithLogger(c, logMock, []string{"2.."})

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
		assert.NoError(t, err)

		resp, err := c.Do(req) //nolint:bodyclose // false positive our request body is nil
		assert.NoError(t, err)
		defer closeHelper(t, resp.Body)
	})

	t.Run("success, response error body", func(t *testing.T) {
		c := &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       0,
		}
		logMock := log.NewMockLogger(gomock.NewController(t))
		logMock.EXPECT().Info(gomock.Any(), "http request", gomock.Any()).Times(1)
		logMock.EXPECT().Error(gomock.Any(), "http response", gomock.Any()).Times(1)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusInternalServerError) }))
		defer srv.Close()

		WithLogger(c, logMock, []string{"2.."})

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
		assert.NoError(t, err)

		resp, err := c.Do(req) //nolint:bodyclose // false positive our request body is nil
		assert.NoError(t, err)
		defer closeHelper(t, resp.Body)
	})

	t.Run("we should not call response log because we are accepting 5xx status", func(t *testing.T) {
		c := &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       0,
		}
		logMock := log.NewMockLogger(gomock.NewController(t))
		logMock.EXPECT().Info(gomock.Any(), "http request", gomock.Any()).Times(1)
		logMock.EXPECT().Info(gomock.Any(), "http response", gomock.Any()).Times(0)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusInternalServerError) }))
		defer srv.Close()

		WithLogger(c, logMock, []string{"[235].."})

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
		assert.NoError(t, err)

		resp, err := c.Do(req) //nolint:bodyclose // false positive our request body is nil
		assert.NoError(t, err)
		defer closeHelper(t, resp.Body)
	})

	t.Run("complete example with trigger, roundtripper logger and accept code", func(t *testing.T) {
		httpClient := NewHTTPClient(&TelemetryMock{})
		t.Setenv("HTTP_TIMEOUT", "10s")

		config := Config{
			Timeout:                0,
			RoundtripperStatusCode: nil,
		}
		err := env.LoadEnv(context.Background(), &config, PrefixHTTP)
		assert.NoError(t, err)

		config.SetTimeout(httpClient)

		logMock := log.NewMockLogger(gomock.NewController(t))
		logMock.EXPECT().Info(gomock.Any(), "http request", gomock.Any()).Times(1)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		defer srv.Close()

		WithLogger(httpClient, logMock, config.RoundtripperStatusCode)

		resp, err := httpClient.Get(srv.URL) //nolint:bodyclose // false positive our request body is nil
		assert.NoError(t, err)
		defer closeHelper(t, resp.Body)
	})
}

func closeHelper(t *testing.T, closer interface{ Close() error }) {
	t.Helper()
	if err := closer.Close(); err != nil {
		t.Error(err)
	}
}
