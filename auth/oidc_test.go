package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	oidc "github.com/coreos/go-oidc"
	"github.com/facily-tech/go-core/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

//nolint:gosec // hardcoded credentials generated from local docker. NOT USED IN ANY ENTERPRISE ENVIRONMENT.
const defaultToken string = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJHYldVYnBHcnllLThoaXFiVWREY0hlSlg4UGE0NVF4NURkVkx0RHFDeHk4In0.eyJleHAiOjE2NjA4NTE4NjgsImlhdCI6MTY2MDg1MTU2OCwianRpIjoiNjI0Y2JjYjQtMWI3Ni00NjhlLTkwZGItYmI3ZTAxZDVjMWRhIiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL3JlYWxtcy9maW5hbmNlIiwiYXVkIjoiZmludGVjaC1sZW5kaW5nLXByb3ZpZGVycyIsInN1YiI6IjNmNTI5NjBhLThhNDItNGNlNi1hMmEyLTA4OWRiNWYwY2U2OCIsInR5cCI6IkJlYXJlciIsImF6cCI6InNlbGxlcnMtYmFja2VuZCIsImFjciI6IjEiLCJhbGxvd2VkLW9yaWdpbnMiOlsiIiwiKiJdLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsidmlldyIsImNoYXJnZSIsImRlZmF1bHQtcm9sZXMtZmluYW5jZSJdfSwicmVzb3VyY2VfYWNjZXNzIjp7InNlbGxlcnMtYmFja2VuZCI6eyJyb2xlcyI6WyJ1bWFfcHJvdGVjdGlvbiJdfX0sInNjb3BlIjoib3BlbmlkIHByb2ZpbGUgZW1haWwgZ29vZC1zZXJ2aWNlIiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJjbGllbnRIb3N0IjoiMTkyLjE2OC4xMjguMSIsImNsaWVudElkIjoic2VsbGVycy1iYWNrZW5kIiwicHJlZmVycmVkX3VzZXJuYW1lIjoic2VydmljZS1hY2NvdW50LXNlbGxlcnMtYmFja2VuZCIsImNsaWVudEFkZHJlc3MiOiIxOTIuMTY4LjEyOC4xIn0.mQBCtFvwdeu91yH-ykqZ3k0cRMUMRkmHytYU1L03W-7yt7NCfPaVm0R_MEMYd8WY_34Joi61LmkMD6KfKjQp01jjIScQHopCkkggAwb3vN4371SeHSdZk5dEGuxi4SXuiosLs-YL7ZmUSAISEd9NrmEsr8UsILpnwomvVtunLx8EtJMVoo8UwuQSE6Y2yikVSTn5T6R-13Z2L3vBl3tnizFL4cRDNhGhn-WqZnND-P5HyVeIwQj3yCKZNnkmppyrfdeN5LYSGJL4uzRqjB1KpysPmcEZFvgXT-EHpYSvMgLOm4uWuG7MynfsSQ_Q6NHaD2L8_cvaeN294vPcG4OUZQ"

type fakeOIDCProvider struct{}

func (o *fakeOIDCProvider) Verifier(config *oidc.Config) *oidc.IDTokenVerifier {
	config.SkipExpiryCheck = true

	return oidc.NewVerifier("http://localhost:8080/realms/finance", &fakeKeySet{}, config)
}

func (o *fakeOIDCProvider) Endpoint() oauth2.Endpoint {
	return oauth2.Endpoint{
		AuthURL:   "",
		TokenURL:  "",
		AuthStyle: oauth2.AuthStyleAutoDetect,
	}
}

type fakeKeySet struct{}

//nolint // copy from oidc parseJWT
func (k *fakeKeySet) VerifySignature(_ context.Context, jwt string) ([]byte, error) {
	parts := strings.Split(jwt, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("oidc: malformed jwt, expected 3 parts got %d", len(parts))
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("oidc: malformed jwt payload: %v", err)
	}

	return payload, nil
}

func TestOIDCAuth(t *testing.T) {
	type fields struct {
		OIDCProvider oidcProvider
		clientID     string
		logger       log.Logger
	}
	type args struct {
		next http.Handler
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "success",
			want: http.StatusOK,
			fields: fields{
				OIDCProvider: &fakeOIDCProvider{},
				clientID:     "fintech-lending-providers",
				logger:       log.NewMockLogger(gomock.NewController(t)),
			},
			args: args{
				next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			},
		},
		{
			name: "fail, wrong client id",
			want: http.StatusUnauthorized,
			fields: fields{
				OIDCProvider: &fakeOIDCProvider{},
				clientID:     "wrong client id",
				logger: func() log.Logger {
					l := log.NewMockLogger(gomock.NewController(t))
					l.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

					return l
				}(),
			},
			args: args{
				next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OIDC{
				oidcProvider: tt.fields.OIDCProvider,
				clientID:     tt.fields.clientID,
				logger:       tt.fields.logger,
			}

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.Header.Add("Authorization", "Bearer "+defaultToken)
			w := httptest.NewRecorder()

			o.Auth(tt.args.next).ServeHTTP(w, r)

			assert.Equal(t, tt.want, w.Code)
		})
	}
}

func TestHasRole(t *testing.T) {
	type args struct {
		ctx  context.Context //nolint:containedctx
		role string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "role found",
			args: args{
				ctx: context.WithValue(context.Background(), contextKey,
					Claims{"realm_access": map[string]interface{}{"roles": []interface{}{"view"}}}),
				role: "view",
			},
			want: true,
		},
		{
			name: "role not found",
			args: args{
				ctx: context.WithValue(context.Background(), contextKey,
					Claims{"realm_access": map[string]interface{}{"roles": []interface{}{"view"}}}),
				role: "charge",
			},
			want: false,
		},
		{
			name: "no key found at ctx",
			args: args{
				ctx:  context.Background(),
				role: "charge",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasRole(tt.args.ctx, tt.args.role); got != tt.want {
				t.Errorf("HasRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasRoleMiddleware(t *testing.T) {
	type args struct {
		role string
	}
	tests := []struct {
		name string
		args args
		req  *http.Request
		want int
	}{
		{
			name: "view role found",
			args: args{
				role: "view",
			},
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)

				return r.WithContext(context.WithValue(context.Background(), contextKey,
					Claims{"realm_access": map[string]interface{}{"roles": []interface{}{"view"}}}))
			}(),
			want: http.StatusOK,
		},
		{
			name: "view role NOT found",
			args: args{
				role: "view",
			},
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)

				return r.WithContext(context.WithValue(context.Background(), contextKey,
					Claims{"realm_access": map[string]interface{}{"roles": []interface{}{"charge"}}}))
			}(),
			want: http.StatusUnauthorized,
		},
		{
			name: "stadard ctx",
			args: args{
				role: "view",
			},
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)

				return r
			}(),
			want: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			HasRoleMiddleware(tt.args.role)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, tt.req)

			assert.Equal(t, tt.want, w.Code)
		})
	}
}

func TestHasScope(t *testing.T) {
	type args struct {
		ctx   context.Context //nolint:containedctx
		scope string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "scope found",
			args: args{
				ctx: context.WithValue(context.Background(), contextKey,
					Claims{"scope": "seller"}),
				scope: "seller",
			},
			want: true,
		},
		{
			name: "scope not found",
			args: args{
				ctx: context.WithValue(context.Background(), contextKey,
					Claims{"scope": "datascience"}),
				scope: "seller",
			},
			want: false,
		},
		{
			name: "no key found at ctx",
			args: args{
				ctx:   context.Background(),
				scope: "charge",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasScope(tt.args.ctx, tt.args.scope); got != tt.want {
				t.Errorf("HasScope() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasScopeMiddleware(t *testing.T) {
	type args struct {
		scope string
	}
	tests := []struct {
		name string
		args args
		req  *http.Request
		want int
	}{
		{
			name: "seller scope found",
			args: args{
				scope: "seller",
			},
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)

				return r.WithContext(context.WithValue(context.Background(), contextKey,
					Claims{"scope": "seller"}))
			}(),
			want: http.StatusOK,
		},
		{
			name: "seller scope NOT found",
			args: args{
				scope: "datascience",
			},
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)

				return r.WithContext(context.WithValue(context.Background(), contextKey,
					Claims{"scope": "seller"}))
			}(),
			want: http.StatusUnauthorized,
		},
		{
			name: "stadard ctx",
			args: args{
				scope: "view",
			},
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)

				return r
			}(),
			want: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			HasScopeMiddleware(tt.args.scope)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, tt.req)

			assert.Equal(t, tt.want, w.Code)
		})
	}
}

func TestGetRootClaim(t *testing.T) {
	type args struct {
		ctx   context.Context //nolint:containedctx
		claim string
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "claim found",
			args: args{
				ctx:   context.WithValue(context.Background(), contextKey, Claims{"clientId": "sellers-backend"}),
				claim: "clientId",
			},
			want: "sellers-backend",
		},
		{
			name: "claim not found",
			args: args{
				ctx:   context.WithValue(context.Background(), contextKey, Claims{"bla": "sellers-backend"}),
				claim: "clientId",
			},
			want: nil,
		},
		{
			name: "claim not found, empty context",
			args: args{
				ctx:   context.Background(),
				claim: "clientId",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRootClaim(tt.args.ctx, tt.args.claim); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRootClaim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkNoAuth(b *testing.B) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, ts.URL, nil)
	if err != nil {
		b.Fatal(err)
	}
	defer ts.Close()

	for i := 0; i < b.N; i++ {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			b.Fatal(err)
		}

		_, err = io.Copy(io.Discard, resp.Body)
		assert.NoError(b, err)

		assert.NoError(b, resp.Body.Close())
	}
}

func BenchmarkAuth(b *testing.B) {
	o := &OIDC{
		oidcProvider: &fakeOIDCProvider{},
		clientID:     "fintech-lending-providers",
		logger:       log.NewMockLogger(gomock.NewController(b)),
	}
	ts := httptest.NewServer(HasRoleMiddleware("view")(o.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))))
	defer ts.Close()

	req, err := http.NewRequestWithContext(context.WithValue(context.Background(), contextKey,
		Claims{"realm_access": map[string]interface{}{"roles": []interface{}{"view"}}}), http.MethodGet, ts.URL, nil)
	if err != nil {
		b.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+defaultToken)

	for i := 0; i < b.N; i++ {
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(b, err)

		_, err = io.Copy(io.Discard, resp.Body)
		assert.NoError(b, err)

		assert.NoError(b, resp.Body.Close())
	}
}
