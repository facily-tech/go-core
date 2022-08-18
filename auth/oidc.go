// Package auth exists to ease the difficult to handle authentication and to some
// extent authorization. auth do not implement the token creation, third-party software
// like keycloak, casdoor etc should be used as issuer.
package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	oidc "github.com/coreos/go-oidc"
	"github.com/facily-tech/go-core/log"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type (
	contextType string
	// Claims of jwt token.
	Claims map[string]interface{}
)

const (
	contextKey contextType = "auth.claims"
	tokenParts int         = 2
)

// OIDC represents our authentication using openid connect and required
// dependencies like logger.
type OIDC struct {
	oidcProvider
	clientID string
	logger   log.Logger
}

// oidcProvider required interface of oidc dependecy.
type oidcProvider interface {
	Verifier(config *oidc.Config) *oidc.IDTokenVerifier
	Endpoint() oauth2.Endpoint
}

// New returns a new OIDC using the provided clientID and issuer to validate
// incoming authentication tokens at OIDC.Auth http middleware method.
func New(log log.Logger, clientID, issuer string) (*OIDC, error) {
	p, err := oidc.NewProvider(context.Background(), issuer)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create oidc provider")
	}

	return &OIDC{
		oidcProvider: p,
		logger:       log,
		clientID:     clientID,
	}, nil
}

// Auth is a middleware used to validate authorization token and populate context
// with stardard claims.
func (o *OIDC) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawAccessToken := r.Header.Get("Authorization")
		if rawAccessToken == "" {
			http.Error(w, "no authorization header", http.StatusBadRequest)
			o.logger.Warn(r.Context(), "auth: empty authorization handler")

			return
		}

		parts := strings.Split(rawAccessToken, " ")
		if len(parts) != tokenParts {
			w.WriteHeader(http.StatusBadRequest)
			o.logger.Warn(r.Context(), "auth: unexpected  authorization size")

			return
		}

		idToken, err := o.Verifier(&oidc.Config{
			ClientID:             o.clientID,
			Now:                  time.Now,
			SupportedSigningAlgs: nil,
			SkipClientIDCheck:    false,
			SkipExpiryCheck:      false,
			SkipIssuerCheck:      false,
		}).Verify(r.Context(), parts[1])
		if err != nil {
			http.Error(w, "cannot validate token: "+err.Error(), http.StatusUnauthorized)
			o.logger.Warn(r.Context(), "auth: invalid token", log.Any("token", parts[1]), log.Error(err))

			return
		}

		resp := Claims{}
		if err := idToken.Claims(&resp); err != nil {
			http.Error(w, "cannot unmarshal claims: "+err.Error(), http.StatusInternalServerError)
			o.logger.Warn(r.Context(), "auth: cannot unmarshal claims", log.Error(err))

			return
		}

		ctx := context.WithValue(r.Context(), contextKey, resp)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// HasRole inspect ctx for role claim, if present returns true. Before using this
// function ctx should be populated with the token claims, generally using
// OIDC.Auth method.
func HasRole(ctx context.Context, role string) bool {
	r := ctx.Value(contextKey)
	if r == nil {
		return false
	}

	claims, ok := r.(Claims)
	if !ok {
		return false
	}
	roles, ok := claims["realm_access"].(map[string][]string)
	if !ok {
		return false
	}

	found := false
	for _, r := range roles["roles"] {
		if role == r {
			found = true

			break
		}
	}

	return found
}

// HasRoleMiddleware wrap HasRole inside an http middleware to prevent access to
// handlers after. OIDC.Auth must be present before this middleware, otherwise
// claims will not be present at http.Request.Context.
func HasRoleMiddleware(role string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !HasRole(r.Context(), role) {
				http.Error(w, "Invalid role", http.StatusUnauthorized)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// HasScope inspect ctx for scope claim, if present returns true. Before using this
// function ctx should be populated with the token claims, generally using
// OIDC.Auth method.
func HasScope(ctx context.Context, scope string) bool {
	s := ctx.Value(contextKey)
	if s == nil {
		return false
	}

	claims, ok := s.(Claims)
	if !ok {
		return false
	}
	scopes, ok := claims["scope"].(string)
	if !ok {
		return false
	}

	found := false
	for _, r := range strings.Split(scopes, " ") {
		if scope == r {
			found = true

			break
		}
	}

	return found
}

// HasScopeMiddleware wrap HasRole inside an http middleware to prevent access to
// handlers after. OIDC.Auth must be present before this middleware, otherwise
// claims will not be present at http.Request.Context.
func HasScopeMiddleware(scope string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !HasScope(r.Context(), scope) {
				http.Error(w, "Invalid scope", http.StatusUnauthorized)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetRootClaim return claim from ctx. ctx should be populated before calling
// GetRootClaim with the token claims, generally using OIDC.Auth method.
func GetRootClaim(ctx context.Context, claim string) interface{} {
	s := ctx.Value(contextKey)
	if s == nil {
		return nil
	}
	claims, ok := s.(Claims)
	if !ok {
		return nil
	}

	return claims[claim]
}
