package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

const (
	bearerPrefix = "Bearer "
)

var (
	bearerSize = len(bearerPrefix)
)

type customerIDContextKey string

// JWTConfig contains the requirements to validate a jwt token.
type JWTConfig struct {
	Secret string `env:"JWT_SECRET,required"`
	Issuer string `env:"JWT_ISSUER,default=https://faci.ly"`
}

type mobileClaim struct {
	Data struct {
		User struct {
			ID string
		}
	}
	jwt.RegisteredClaims
}

// JWTMW middleware checks for jwt token and if present validate and populate
// with custom claim ID. ID can be retrieved using GetCustomerID.
func (j JWTConfig) JWTMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodedToken := getToken(r.Header.Get("authorization"))
		if encodedToken == "" {
			http.Error(w, "missing authorization token", http.StatusUnauthorized)

			return
		}

		token, err := jwt.ParseWithClaims(encodedToken, &mobileClaim{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}

			iss, err := t.Claims.GetIssuer()
			if err != nil {
				return nil, errors.WithStack(err)
			}
			if iss != j.Issuer {
				return nil, jwt.ErrTokenInvalidIssuer
			}

			return []byte(j.Secret), nil
		})
		if err != nil {
			http.Error(w, "{'message': 'token isn't valid: "+ err.Error() + "'}", http.StatusUnauthorized)

			return
		}

		claims, ok := token.Claims.(*mobileClaim)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		customerID, err := strconv.ParseInt(claims.Data.User.ID, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		r = r.WithContext(SetCustomerID(r.Context(), int(customerID)))

		next.ServeHTTP(w, r)
	})
}

// GetCustomerID retrieve id from ctx. Return -1 of type assertion fail.
func GetCustomerID(ctx context.Context) int {
	id, ok := ctx.Value(customerIDContextKey("customerID")).(int)
	if !ok {
		return -1
	}

	return id
}

// SetCustomerID returns a new context using ctx as parent and inserting id.
func SetCustomerID(ctx context.Context, id int) context.Context {
	return context.WithValue(ctx, customerIDContextKey("customerID"), id)
}

func getToken(authorizationHeader string) string {
	if !strings.HasPrefix(authorizationHeader, bearerPrefix) {
		return ""
	}

	return authorizationHeader[bearerSize:]
}
