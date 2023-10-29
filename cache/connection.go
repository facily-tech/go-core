package cache

import (
	"context"
	"strings"
	"time"

	"github.com/facily-tech/go-core/env"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/redis/go-redis.v9"
)

// CachePrefix is the prefix for the environment variables.
const CachePrefix = "CACHE_"

// ClientI is the interface for the cache.
var _ ClientI = (*Client)(nil)

// ErrKeyMiss is the error returned when the key does not exist.
var ErrKeyMiss = errors.New("key does not exist")

// Client is the interface for the cache.
type Client struct {
	client redis.UniversalClient
}

type config struct {
	Address              string `env:"ADDR,required"`
	Password             string `env:"PASSWORD"`
	DB                   int    `env:"DB,default=10"`
	TracerDatadogEnabled bool   `env:"DD_TRACE,default=true"`
}

// InitCache initializes the cache.
func InitCache() (*Client, error) {
	return initCache(CachePrefix)
}

func initCache(cachePrefix string) (*Client, error) {
	cacheConfig, err := loadEnv(cachePrefix)
	if err != nil {
		return nil, err
	}

	return openConn(cacheConfig)
}

func openConn(cacheConfig *config) (*Client, error) {
	var rdb redis.UniversalClient

	opts := &redis.UniversalOptions{
		Addrs:    strings.Split(cacheConfig.Address, ","),
		Password: cacheConfig.Password,
		DB:       cacheConfig.DB,
	}

	if cacheConfig.TracerDatadogEnabled {
		rdb = redistrace.NewClient(opts.Simple())
	} else {
		rdb = redis.NewUniversalClient(opts)
	}

	timeout, c := context.WithTimeout(context.Background(), time.Minute)
	defer c()

	if err := rdb.Ping(timeout).Err(); err != nil {
		return nil, errors.Wrap(err, "cannot ping redis")
	}

	return &Client{rdb}, nil
}

// loadEnv loads the environment variables.
func loadEnv(cachePrefix string) (*config, error) {
	var cacheConfig config
	if err := env.LoadEnv(context.Background(), &cacheConfig, cachePrefix); err != nil {
		return nil, errors.Wrap(err, "cannot load db environment variable")
	}

	return &cacheConfig, nil
}

// Set define the value of a key.
func (r *Client) Set(ctx context.Context, key string, value string, expire time.Duration) error {
	return errors.Wrapf(r.client.Set(ctx, key, value, expire).Err(), "cannot set key '%s'", key)
}

// Get gets the value of a key.
func (r *Client) Get(ctx context.Context, key string) (string, error) {
	v, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", errors.Wrapf(ErrKeyMiss, "miss trying to get key '%s'", key)
	}
	if err != nil {
		return "", errors.Wrapf(err, "cannot get key '%s'", key)
	}

	return v, nil
}

// Del deletes one or more keys.
func (r *Client) Del(ctx context.Context, keys ...string) error {
	return errors.Wrap(r.client.Del(ctx, keys...).Err(), "cannot delete")
}
