package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/facily-tech/go-core/log"
)

const (
	expectedAuthSize = 2
)

type OAuth2Fixed struct {
	Logger log.Logger
	Token  string
}

func (o OAuth2Fixed) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			o.Logger.Warn(ctx, "authorization header not found")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		typeCred := strings.Split(authHeader, " ")
		if len(typeCred) != expectedAuthSize {
			o.Logger.Warn(ctx, "unexpected authorization size", log.Any("authorization", typeCred))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if typeCred[1] != o.Token {
			o.Logger.Warn(ctx, "token didn't match", log.Any("token", typeCred))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (o OAuth2Fixed) AccessToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Content-type", "application/json")
	fmt.Fprintf(w, `
{
  "access_token":"%s",
  "token_type":"bearer",
  "expires_in":3600
}`, o.Token)
}
