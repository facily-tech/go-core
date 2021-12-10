package http

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

const BasicAuthConfigPrefix = "BASIC_AUTH_"

type BasicAuth struct {
	Realm        string            `env:"REALM,default=Oracle Retail Publishing"`
	UserPassword map[string]string `env:"USER_PASSWORD,required"`
}

func (ba BasicAuth) Middleware(next http.Handler) http.Handler {
	return middleware.BasicAuth(ba.Realm, ba.UserPassword)(next)
}
