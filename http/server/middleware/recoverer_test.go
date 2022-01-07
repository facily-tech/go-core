package middleware

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/facily-tech/go-core/log"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func testRequest(ctx context.Context, t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	t.Helper()
	req, err := http.NewRequestWithContext(ctx, method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)

		return nil, ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)

		return nil, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)

		return nil, ""
	}

	return resp, string(respBody)
}

func panicingHandler(http.ResponseWriter, *http.Request) {
	panic("foo")
}

func TestRecoverer_withPanic(t *testing.T) {
	ctx := context.Background()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogger := log.NewMockLogger(mockCtrl)
	mockLogger.EXPECT().Error(
		gomock.Not(gomock.Nil()),
		panicErrorRecovered,
		gomock.Not(gomock.Nil()),
	).Times(1)

	r := chi.NewRouter()
	r.Use(Recoverer(mockLogger))
	r.Get("/", panicingHandler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	res, _ := testRequest(ctx, t, ts, "GET", "/", nil)
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Error(err)
		}
	}()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestRecoverer_withOutPanic(t *testing.T) {
	ctx := context.Background()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLogger := log.NewMockLogger(mockCtrl)
	mockLogger.EXPECT().Error(
		gomock.Not(gomock.Nil()),
		panicErrorRecovered,
		gomock.Not(gomock.Nil()),
	).Times(0)

	r := chi.NewRouter()
	r.Use(Recoverer(mockLogger))
	r.Get("/", func(http.ResponseWriter, *http.Request) {})

	ts := httptest.NewServer(r)
	defer ts.Close()

	res, _ := testRequest(ctx, t, ts, "GET", "/", nil)
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.Error(err)
		}
	}()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}
