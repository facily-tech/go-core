//nolint:revive,stylecheck // yeah mixpanel use bad names
package analytics

import (
	"testing"
	"time"

	"github.com/dukex/mixpanel"
	"github.com/facily-tech/go-core/log"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type MixMock struct {
	// Create a mixpanel event using the track api
	trackF func(distinctId string, eventName string, e *mixpanel.Event) error
	// Create a mixpanel event using the import api
	importF func(distinctId string, eventName string, e *mixpanel.Event) error
	// Set properties for a mixpanel user.
	// Deprecated: Use UpdateUser instead
	updateF func(distinctId string, u *mixpanel.Update) error
	// Set properties for a mixpanel user.
	updateUserF func(distinctId string, u *mixpanel.Update) error
	// Set properties for a mixpanel group.
	updateGroupF func(groupKey string, groupId string, u *mixpanel.Update) error
	// Create an alias for an existing distinct id
	aliasF func(distinctId string, newId string) error
}

// Create a mixpanel event using the track api.
func (m *MixMock) Track(distinctId, eventName string, e *mixpanel.Event) error {
	return m.trackF(distinctId, eventName, e)
}

// Create a mixpanel event using the import api.
func (m *MixMock) Import(distinctId, eventName string, e *mixpanel.Event) error {
	return m.importF(distinctId, eventName, e)
}

// Set properties for a mixpanel user.
// Deprecated: Use UpdateUser instead.
func (m *MixMock) Update(distinctId string, u *mixpanel.Update) error {
	return m.updateF(distinctId, u)
}

// Set properties for a mixpanel user.
func (m *MixMock) UpdateUser(distinctId string, u *mixpanel.Update) error {
	return m.updateUserF(distinctId, u)
}

// Set properties for a mixpanel group.
func (m *MixMock) UpdateGroup(groupKey, groupId string, u *mixpanel.Update) error {
	return m.updateGroupF(groupKey, groupId, u)
}

// Create an alias for an existing distinct id.
func (m *MixMock) Alias(distinctId, newId string) error {
	return m.aliasF(distinctId, newId)
}

func TestTrackSync(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		a := Analytics{
			MixpanelEvent: nil,
			Logger:        nil,
			token:         "",
		}
		a.WithMixpanelURL("token", "")
		logger := log.NewMockLogger(gomock.NewController(t))
		a.Logger = logger
		a.MixpanelEvent = &mixpanel.Mock{People: make(map[string]*mixpanel.MockPeople)}

		err := a.TrackSync("eventName", map[string]interface{}{"propertie": "value"})
		assert.NoError(t, err)
	})
	t.Run("no token", func(t *testing.T) {
		a := Analytics{
			MixpanelEvent: nil,
			Logger:        nil,
			token:         "",
		}
		logger := log.NewMockLogger(gomock.NewController(t))
		a.Logger = logger
		a.MixpanelEvent = &MixMock{}

		err := a.TrackSync("eventName", map[string]interface{}{"propertie": "value"})
		assert.Error(t, err)
	})
	t.Run("mixpanel event must fail", func(t *testing.T) {
		a := Analytics{
			MixpanelEvent: nil,
			Logger:        nil,
			token:         "token",
		}
		logger := log.NewMockLogger(gomock.NewController(t))
		a.Logger = logger
		//nolint:exhaustruct // accept default values at structs
		a.MixpanelEvent = &MixMock{
			trackF: func(distinctId, eventName string, e *mixpanel.Event) error {
				return errors.New("random error")
			},
		}

		err := a.TrackSync("eventName", map[string]interface{}{"propertie": "value"})
		assert.Error(t, err)
	})
	t.Run("async track", func(t *testing.T) {
		a := Analytics{
			MixpanelEvent: nil,
			Logger:        nil,
			token:         "token",
		}
		logger := log.NewMockLogger(gomock.NewController(t))
		a.Logger = logger

		done := make(chan struct{})
		a.WithWorkerPool(poolSize)
		//nolint:exhaustruct // accept default values at structs
		a.MixpanelEvent = &MixMock{
			trackF: func(distinctId, eventName string, e *mixpanel.Event) error {
				defer func() { done <- struct{}{} }()

				return nil
			},
		}

		a.Track("eventName", map[string]interface{}{"propertie": "value"})
		select {
		case <-time.After(time.Minute):
			t.Error("Track timeout")
		case <-done:
		}
	})
}
