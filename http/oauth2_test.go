package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/facily-tech/go-core/log"
	"github.com/stretchr/testify/assert"
)

type mockedLogger struct {
	Message string
	Fields  []log.Field
}

// Error logs a message at ErrorLevel. The message includes any fields passed at the log site, as well as any
// fields accumulated on the logger.
func (l *mockedLogger) Error(ctx context.Context, msg string, fields ...log.Field) {}

// Debug logs a message at DebugLevel. The message includes any fields passed at the log site, as well as any
// fields accumulated on the logger.
func (l *mockedLogger) Debug(ctx context.Context, msg string, fields ...log.Field) {}

// Fatal logs a message at FatalLevel. The message includes any fields passed at the log site, as well as any
// fields accumulated on the logger. The logger then calls os.Exit(1), even if logging at FatalLevel is disabled.
// Defer aren't executed before exit! Use only in appropriated places like simple main() without defer.
func (l *mockedLogger) Fatal(ctx context.Context, msg string, fields ...log.Field) {}

// Info logs a message at InfoLevel. The message includes any fields passed at the log site, as well as any fields
// accumulated on the logger.
func (l *mockedLogger) Info(ctx context.Context, msg string, fields ...log.Field) {}

// Panic logs a message at PanicLevel. The message includes any fields passed at the log site, as well as any fields
// accumulated on the logger. The logger then panics, even if logging at PanicLevel is disabled.
func (l *mockedLogger) Panic(ctx context.Context, msg string, fields ...log.Field) {}

// Warn logs a message at WarnLevel. The message includes any fields passed at the log site, as well as any fields
// accumulated on the logger.
func (l *mockedLogger) Warn(ctx context.Context, msg string, fields ...log.Field) {
	l.Message = msg
	l.Fields = fields
}

type mockedHandler struct {
	Calls uint
}

func (h *mockedHandler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {
	h.Calls++
}

func TestOAuth2Fixed_Middleware(t *testing.T) {
	type fields struct {
		Logger log.Logger
		Token  string
	}
	type args struct {
		next http.Handler
	}
	type want struct {
		msg       string
		fields    []log.Field
		nextCalls uint
		status    int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "success",
			fields: fields{
				Logger: &mockedLogger{},
				Token:  "1",
			},
			args: args{
				next: &mockedHandler{},
			},
			want: want{
				msg:       "",
				fields:    nil,
				nextCalls: 1,
				status:    http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := OAuth2Fixed{
				Logger: tt.fields.Logger,
				Token:  tt.fields.Token,
			}
			got := o.Middleware(tt.args.next)
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", "Bearer 1")
			got.ServeHTTP(rr, req)

			logger, ok := tt.fields.Logger.(*mockedLogger)
			assert.True(t, ok)
			handler, ok := tt.args.next.(*mockedHandler)
			assert.True(t, ok)

			assert.Equal(t, tt.want.fields, logger.Fields)
			assert.Equal(t, tt.want.msg, logger.Message)
			assert.Equal(t, tt.want.nextCalls, handler.Calls)
			assert.Equal(t, tt.want.status, rr.Code)
		})
	}
}
