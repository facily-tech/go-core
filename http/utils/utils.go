package http

import (
	"context"
	"net/http"

	"github.com/facily-tech/go-core/log"
)

type ModeEnum int

const (
	ServerMode ModeEnum = iota
	ClientMode
)

func StatusLevel(logger log.Logger, status int, mode ModeEnum) func(ctx context.Context, msg string, fields ...log.Field) {
	switch {
	case status <= 0:
		return logger.Warn
	case status < http.StatusBadRequest: // for codes in 100s, 200s, 300s
		if mode == ClientMode && status >= http.StatusMultipleChoices {
			return logger.Warn
		}
		return logger.Info
	case status >= http.StatusBadRequest && status < http.StatusInternalServerError:
		if mode == ClientMode {
			return logger.Error
		}
		return logger.Warn
	case status >= http.StatusInternalServerError:
		return logger.Error
	default:
		return logger.Info
	}
}
