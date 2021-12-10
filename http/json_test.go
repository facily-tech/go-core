package http

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/facily-tech/go-core/log"
	"github.com/stretchr/testify/assert"
)

func TestToJSON(t *testing.T) {
	ctx := context.Background()
	zap, err := log.NewLoggerZap(
		log.ZapConfig{},
	)
	assert.NoError(t, err)

	type args struct {
		logger log.Logger
		input  interface{}
		w      *httptest.ResponseRecorder
		status int
	}
	tests := []struct {
		name     string
		args     args
		want     string
		wantCode int
	}{
		{
			name: "success",
			args: args{
				logger: zap,
				input:  map[string]string{"field": "value"},
				w:      httptest.NewRecorder(),
				status: http.StatusOK,
			},
			want:     `{"field":"value"}` + "\n",
			wantCode: http.StatusOK,
		},
		{
			name: "fail on complex types",
			args: args{
				logger: zap,
				input:  map[string]func(){"field": func() {}},
				w:      httptest.NewRecorder(),
				status: http.StatusOK,
			},
			want:     ``,
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ToJSON(ctx, tt.args.logger, tt.args.input, tt.args.w, tt.args.status)
			r := tt.args.w.Result()
			assert.Equal(t, tt.want, tt.args.w.Body.String())
			assert.Equal(t, tt.wantCode, r.StatusCode)
			r.Body.Close()
		})
	}
}

type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) {
	return 0, os.ErrClosed
}

func TestFromJSON(t *testing.T) {
	type outtesting struct {
		Key string
	}
	zap, err := log.NewLoggerZap(log.ZapConfig{})
	assert.NoError(t, err)

	type args struct {
		logger log.Logger
		output interface{}
		r      io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    outtesting
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				logger: zap,
				output: &outtesting{},
				r:      strings.NewReader(`{"key": "value"}`),
			},
			want:    outtesting{Key: "value"},
			wantErr: false,
		},
		{
			name: "fail, error while reading io.Reader",
			args: args{
				logger: zap,
				output: &outtesting{},
				r:      errorReader{},
			},
			want:    outtesting{Key: ""},
			wantErr: true,
		}, {
			name: "fail, invalid json",
			args: args{
				logger: zap,
				output: &outtesting{},
				r:      strings.NewReader(`{key": "value}`),
			},
			want:    outtesting{Key: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := FromJSON(tt.args.logger, tt.args.output, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("FromJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, &tt.want, tt.args.output)
		})
	}
}
