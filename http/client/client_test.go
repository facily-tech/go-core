/*
Package client make easy to use an http client
*/
package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/facily-tech/go-core/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestWithLogger(t *testing.T) {
	t.Run("success, request info", func(t *testing.T) {
		c := &http.Client{}
		logMock := log.NewMockLogger(gomock.NewController(t))
		logMock.EXPECT().Info(gomock.Any(), "http request", gomock.Any()).Times(1)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		defer srv.Close()

		WithLogger(c, logMock)(c)

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

		WithLogger(c, logMock)(c)

		resp, err := c.Get(srv.URL)
		assert.NoError(t, err)
		defer resp.Body.Close()

	})
}
