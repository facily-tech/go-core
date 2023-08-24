package middleware

import (
	"bytes"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	utils "github.com/facily-tech/go-core/http/utils"
	"github.com/facily-tech/go-core/log"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

// DontLogBodyOnSuccess default prefix to not log response in case of success.
var DontLogBodyOnSuccess = []string{
	"/swagger",
	"/metrics",
	"/health",
}

// ContentTypeToLogBody is a list of content-type that the body will be logged.
var ContentTypeToLogBody = []string{
	"application/json",
	"application/xml",
	"text/json",
	"text/xml",
	"text/plain",
}

// wrapWriter implements http.ResponseWriter and saves status and size for logging.
type wrapWriter struct {
	http.ResponseWriter
	status      int
	size        int
	wroteHeader bool
	bodyBuffer  bytes.Buffer
}

// Write writes the data to the connection as part of an HTTP reply.
//
// If WriteHeader has not yet been called, Write calls
// WriteHeader(http.StatusOK) before writing the data. If the Header
// does not contain a Content-Type line, Write adds a Content-Type set
// to the result of passing the initial 512 bytes of written data to
// DetectContentType. Additionally, if the total size of all written
// data is under a few KB and there are no Flush calls, the
// Content-Length header is added automatically.
//
// Depending on the HTTP protocol version and the client, calling
// Write or WriteHeader may prevent future reads on the
// Request.Body. For HTTP/1.x requests, handlers should read any
// needed request body data before writing the response. Once the
// headers have been flushed (due to either an explicit Flusher.Flush
// call or writing enough data to trigger a flush), the request body
// may be unavailable. For HTTP/2 requests, the Go HTTP server permits
// handlers to continue to read the request body while concurrently
// writing the response. However, such behavior may not be supported
// by all HTTP/2 clients. Handlers should read before writing if
// possible to maximize compatibility.
func (r *wrapWriter) Write(data []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	bn, err := r.bodyBuffer.Write(data)
	if err != nil {
		return bn, errors.WithStack(err)
	}

	n, err := r.ResponseWriter.Write(data)
	r.size += n

	return n, errors.Wrap(err, "log wrapper unable to write response")
}

// WriteHeader sends an HTTP response header with the provided
// status code.
//
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
//
// The provided code must be a valid HTTP 1xx-5xx status code.
// Only one header may be written. Go does not currently
// support sending user-defined 1xx informational headers,
// with the exception of 100-continue response header that the
// Server sends automatically when the Request.Body is read.
func (r *wrapWriter) WriteHeader(statusCode int) {
	if r.wroteHeader {
		return
	}
	r.status = statusCode
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(statusCode)
}

// Logger middleware is a middleware to log everything request that was receved by API.
func Logger(logger log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writer := &wrapWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
				size:           0,
				wroteHeader:    false,
				bodyBuffer:     bytes.Buffer{},
			}
			tt := time.Now()

			authBkp := make([]string, len(r.Header["Authorization"]))
			if v, exist := r.Header["Authorization"]; exist {
				copy(authBkp, v)
				r.Header["Authorization"] = []string{"****"}
			}

			reqbody, err := httputil.DumpRequest(r, true)
			if err != nil {
				logger.Error(r.Context(), "can't open request", log.Error(err))
			}

			if _, exist := r.Header["Authorization"]; exist {
				copy(r.Header["Authorization"], authBkp)
			}

			logger.Info(r.Context(), "request",
				log.Any("method", r.Method),
				log.Any("path", r.URL.Path),
				log.Any("from", getIP(r)),
				log.Any("body", getBodyAsString(r.Header.Get("Content-Type"), reqbody)),
			)

			next.ServeHTTP(writer, r)

			for _, v := range DontLogBodyOnSuccess {
				if strings.HasPrefix(r.URL.Path, v) && writer.status < http.StatusBadRequest {
					return
				}
			}

			resBodyLog := getBodyAsString(
				writer.ResponseWriter.Header().Get("Content-Type"),
				writer.bodyBuffer.Bytes(),
			)

			utils.StatusLevel(logger, writer.status, utils.ServerMode)(
				r.Context(), "response",
				log.Any("method", r.Method),
				log.Any("path", r.URL.Path),
				log.Any("from", getIP(r)),
				log.Any("status", writer.status),
				log.Any("size_bytes", writer.size),
				log.Any("elapsed_seconds", time.Since(tt).Seconds()),
				log.Any("elapsed", time.Since(tt).String()),
				log.Any("body", resBodyLog),
			)
		})
	}
}

// getIP returns the IP of the request.
func getIP(r *http.Request) string {
	IP := r.RemoteAddr
	// try to get ip from reverse proxy header
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		IP = ip
	}

	return IP
}

// getBodyAsString returns the body as string if the content-type is in the list of ContentTypeToLogBody.
func getBodyAsString(contentTpe string, body []byte) string {
	if len(contentTpe) == 0 || slices.Index(ContentTypeToLogBody, contentTpe) > 0 {
		return string(body)
	}

	return "***"
}
