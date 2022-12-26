package analytics

import (
	"context"

	"github.com/dukex/mixpanel"
	"github.com/facily-tech/go-core/env"
	"github.com/facily-tech/go-core/log"
	"github.com/gammazero/workerpool"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	prefix   = "ANALYTICS_"
	poolSize = 5
)

type config struct {
	Token string `env:"MIXPANEL_TOKEN"`
}

// Analytics is an wrapper over mixpanel package.
type Analytics struct {
	MixpanelEvent mixpanel.Mixpanel
	Logger        log.Logger

	token string
	wp    workerPool
}

type workerPool interface {
	StopWait()
	Submit(func())
}

var (
	// DefaultAnalytics is used if you didn't create your own instance.
	DefaultAnalytics = &Analytics{}

	// ErrUnexpectedType error when we can't cast/type assert interface to (T).
	ErrUnexpctedType = errors.New("cannot cast to desired type")
	// ErrEmptyToken error when token is empty.
	ErrEmptyToken = errors.New("token cannot by empty, use WithMixpanelURL")
)

//nolint:gochecknoinits // try to load token from env.
func init() {
	var c config
	if err := env.LoadEnv(context.Background(), &c, prefix); err != nil {
		panic(err)
	}

	logger, err := log.NewLoggerZap(log.ZapConfig{})
	if err != nil {
		panic(err)
	}

	DefaultAnalytics.MixpanelEvent = mixpanel.New(c.Token, "")
	DefaultAnalytics.token = c.Token
	DefaultAnalytics.Logger = logger
	DefaultAnalytics.wp = workerpool.New(poolSize)
}

// WithWorkerPool stop current work pool if already running and creates
// a new one with size workers.

// WithWorkerPool is a wrapper around DefaultClient.WithWorkerPool.
func WithWorkerPool(size int) {
	if DefaultAnalytics.wp != nil {
		DefaultAnalytics.wp.StopWait()
	}
	DefaultAnalytics.wp = workerpool.New(size)
}

// WithWorkerPool stop current work pool if already running and creates
// a new one with size workers.
func (a *Analytics) WithWorkerPool(size int) {
	if a.wp != nil {
		a.wp.StopWait()
	}
	a.wp = workerpool.New(size)
}

// WithLogger change current log.Logger to externalLogger.

// WithLogger is a wrapper around DefaultClient.WithLogger.
func WithLogger(externalLogger log.Logger) {
	DefaultAnalytics.WithLogger(externalLogger)
}

// WithLogger change current log.Logger to externalLogger.
func (a *Analytics) WithLogger(externalLogger log.Logger) {
	a.Logger = externalLogger
}

// WithMixpanelURL client token and url. url can be empty and will default to
// "https://api.mixpanel.com".

// WithMixpanelURL is a wrapper around DefaultClient.WithMixpanelURL.
func WithMixpanelURL(token, url string) {
	DefaultAnalytics.MixpanelEvent = mixpanel.New(token, url)
}

// WithMixpanelURL client token and url. url can be empty and will default to
// "https://api.mixpanel.com".
func (a *Analytics) WithMixpanelURL(token, url string) {
	a.token = token
	a.MixpanelEvent = mixpanel.New(token, url)
}

// Track queue an eventName with the following properties to be sent.
func (a *Analytics) Track(eventName string, properties map[string]interface{}) {
	a.wp.Submit(func() {
		err := a.dispatchEvent(eventName, properties)
		if err != nil {
			a.Logger.Error(context.Background(), "Error sending event to mixpanel", log.Error(err))
		}
	})
}

// Track queue an eventName with the following properties to be sent.

// Track is a wrapper around DefaultClient.Track.
func Track(eventName string, properties map[string]interface{}) {
	DefaultAnalytics.Track(eventName, properties)
}

// TrackSync send an eventName with the following properties.

// TrackSync is a wrapper around DefaultClient.TrackSync.
func TrackSync(eventName string, properties map[string]interface{}) error {
	return DefaultAnalytics.dispatchEvent(eventName, properties)
}

// TrackSync send an eventName with the following properties.
func (a *Analytics) TrackSync(eventName string, properties map[string]interface{}) error {
	return a.dispatchEvent(eventName, properties)
}

func (a *Analytics) dispatchEvent(eventName string, properties map[string]interface{}) error {
	distinctID := uuid.NewString()
	if properties["distinctId"] != nil {
		d, ok := properties["distinctId"].(string)
		if !ok {
			return errors.Wrap(ErrUnexpctedType, "distinctId is not a string")
		}
		distinctID = d
	}

	if a.token == "" {
		return errors.WithStack(ErrEmptyToken)
	}

	if err := a.MixpanelEvent.Track(distinctID, eventName, &mixpanel.Event{
		IP:         "",
		Timestamp:  nil,
		Properties: properties,
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Close wait for events to be sent and stop worker pool.
func (a *Analytics) Close() {
	a.wp.StopWait()
}

// Close wait for events to be sent and stop worker pool.

// Close is a wrapper around DefaultClient.Close.
func Close() {
	DefaultAnalytics.wp.StopWait()
}
