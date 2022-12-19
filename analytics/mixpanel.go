package analytics

import (
	"context"

	"github.com/dukex/mixpanel"
	"github.com/facily-tech/go-core/env"
	"github.com/facily-tech/go-core/log"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const prefix = "ANALYTICS_"

type config struct {
	Token string `env:"MIXPANEL_TOKEN,required"`
}

type Analytics struct {
	MixpanelEvent mixpanel.Mixpanel
	Logger        log.Logger
}

var DefaultAnalytics = &Analytics{}
var UnexpctedType = errors.New("cannot cast to desired type")

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
	DefaultAnalytics.Logger = logger
}

func WithLogger(externalLogger log.Logger) {
	DefaultAnalytics.WithLogger(externalLogger)
}

func (a *Analytics) WithLogger(externalLogger log.Logger) {
	a.Logger = externalLogger
}

func WithMixpanelURL(token, url string) {
	DefaultAnalytics.MixpanelEvent = mixpanel.New(token, url)
}

func (a *Analytics) WithMixpanelURL(token, url string) {
	a.MixpanelEvent = mixpanel.New(token, url)
}

func (a *Analytics) Track(eventName string, properties map[string]interface{}) {
	go func(eventName string, properties map[string]interface{}) {
		err := a.dispatchEvent(eventName, properties)
		if err != nil {
			a.Logger.Error(context.Background(), "Error sending event to mixpanel", log.Error(err))
		}
	}(eventName, properties)
}

func Track(eventName string, properties map[string]interface{}) {
	DefaultAnalytics.Track(eventName, properties)
}

func TrackSync(eventName string, properties map[string]interface{}) error {
	return DefaultAnalytics.dispatchEvent(eventName, properties)
}

func (a *Analytics) TrackSync(eventName string, properties map[string]interface{}) error {
	return a.dispatchEvent(eventName, properties)
}

func (a *Analytics) dispatchEvent(eventName string, properties map[string]interface{}) error {
	distinctID := uuid.NewString()
	if properties["distinctId"] != nil {
		d, ok := properties["distinctId"].(string)
		if !ok {
			return errors.Wrap(UnexpctedType, "distinctId is not a string")
		}
		distinctID = d
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
