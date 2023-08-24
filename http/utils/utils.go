package utils

import (
	"context"
	"net/http"

	"github.com/facily-tech/go-core/log"
)

// ModeEnum is the identifier of the http mode.
type ModeEnum int

const (
	// ServerMode is when the http is used as a server.
	ServerMode ModeEnum = iota + 1
	// ClientMode is when the http is used as a client.
	ClientMode
)

// StatusLevel is a function that return the log level based on the status code.
func StatusLevel(logger log.Logger, status int, mode ModeEnum) func(ctx context.Context,
	msg string, fields ...log.Field) {
	switch {
	case status <= 0:
		return logger.Warn
	case status < http.StatusBadRequest: // for codes in 100s, 200s, 300s
		if mode == ClientMode && status >= http.StatusMultipleChoices {
			return logger.Warn
		}

		return logger.Info
	case status >= http.StatusBadRequest && status < http.StatusInternalServerError:
		if shouldLogError(status, mode) {
			return logger.Error
		}

		return logger.Warn
	case status >= http.StatusInternalServerError:
		return logger.Error
	default:
		return logger.Info
	}
}

func shouldLogError(status int, mode ModeEnum) bool {
	return mode == ClientMode && status != http.StatusNotFound
}
