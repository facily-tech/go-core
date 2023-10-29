package cache

import (
	"context"
	"time"
)

// ClientI is the interface for the cache.
type ClientI interface {
	Set(context.Context, string, string, time.Duration) error
	Get(context.Context, string) (string, error)
	Del(context.Context, ...string) error
}
