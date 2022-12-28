package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/facily-tech/go-core/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	t.Run("success simplest GET", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		writer := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		mockLog := log.NewMockLogger(gomock.NewController(t))
		mockLog.EXPECT().Info(
			gomock.Any(), // context
			"request",
			log.Any("method", http.MethodGet),
			log.Any("path", "/"),
			gomock.Any(), // from
			log.Any("body", "GET / HTTP/1.1\r\nHost: example.com\r\n\r\n"),
		)

		mockLog.EXPECT().Info(
			gomock.Any(), // context
			"response",
			log.Any("method", http.MethodGet),
			log.Any("path", "/"),
			gomock.Any(), // from
			log.Any("status", http.StatusOK),
			log.Any("size_bytes", 0),
			gomock.Any(), // elapsed_seconds
			gomock.Any(), // elapsed
			log.Any("body", ""),
		)

		Logger(mockLog)(handler).ServeHTTP(writer, req)
		response := writer.Result()
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.NoError(t, response.Body.Close())
	})

	t.Run("success GET with Authorization header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add("Authorization", "secret")

		writer := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "secret", r.Header.Get("Authorization"))
		})
		mockLog := log.NewMockLogger(gomock.NewController(t))
		mockLog.EXPECT().Info(
			gomock.Any(), // context
			"request",
			log.Any("method", http.MethodGet),
			log.Any("path", "/"),
			gomock.Any(), // from
			log.Any("body", "GET / HTTP/1.1\r\nHost: example.com\r\nAuthorization: ****\r\n\r\n"),
		)

		mockLog.EXPECT().Info(
			gomock.Any(), // context
			"response",
			log.Any("method", http.MethodGet),
			log.Any("path", "/"),
			gomock.Any(), // from
			log.Any("status", http.StatusOK),
			log.Any("size_bytes", 0),
			gomock.Any(), // elapsed_seconds
			gomock.Any(), // elapsed
			log.Any("body", ""),
		)

		Logger(mockLog)(handler).ServeHTTP(writer, req)
		response := writer.Result()
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.NoError(t, response.Body.Close())
	})
}
