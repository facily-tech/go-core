package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

//nolint:bodyclose // false positive, body is bytes.Buffer
func TestJWTConfig_JWTMW(t *testing.T) {
	type fields struct {
		Secret string
		Issuer string
	}
	type args struct {
		next http.Handler
	}
	tests := []struct {
		args     args
		fields   fields
		name     string
		wantCode int
	}{
		{
			name: "every input is valid, no ctx validation expect success",
			fields: fields{
				Secret: "fake secret",
				Issuer: "https://faci.ly",
			},
			args: args{
				next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			},
			wantCode: http.StatusOK,
		},
		{
			name: "validated token, expect success",
			fields: fields{
				Secret: "fake secret",
				Issuer: "https://faci.ly",
			},
			args: args{
				next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					id := GetCustomerID(r.Context())
					assert.Equal(t, 5811308, id)
				}),
			},
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := JWTConfig{
				Secret: tt.fields.Secret,
				Issuer: tt.fields.Issuer,
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.Header.Add("authorization", "Bearer "+"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwczovL2ZhY2kubHkiLCJpYXQiOjE2NzU0MDAzNDgsIm5iZiI6MTY3NTQwMDM0OCwiZXhwIjoxNjgzNzgxMTQ4LCJkYXRhIjp7InVzZXIiOnsiaWQiOiI1ODExMzA4In19fQ.e8HoV5Oj-WbVhSsrqiHW8lJWB5XWhIiXCk4GxpU06hs")

			j.JWTMW(tt.args.next).ServeHTTP(w, r)

			assert.Equal(t, tt.wantCode, w.Result().StatusCode)
			assert.Equal(t, "", w.Body.String())
		})
	}
}

func TestGetCustomerID(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		args args
		name string
		want int
	}{
		{
			name: "set valid ctx value, expect success",
			args: args{
				ctx: SetCustomerID(context.Background(), 123456),
			},
			want: 123456,
		},
		{
			name: "set invalid ctx value, expect -1",
			args: args{
				ctx: context.WithValue(context.Background(), customerIDContextKey("customerID"), "a"),
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCustomerID(tt.args.ctx); got != tt.want {
				t.Errorf("GetCustomerID() = %v, want %v", got, tt.want)
			}
		})
	}
}
