package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/facily-tech/go-core/log"
)

const panicErrorRecovered string = "panic recovered on middleware recoverer"

func Recoverer(logger log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				rvr := recover()
				logger.Error(
					r.Context(),
					panicErrorRecovered,
					log.Any("recover", rvr),
					log.Any("debug", string(debug.Stack())),
				)
				w.WriteHeader(http.StatusInternalServerError)
			}()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
