package rules

import (
	"geoip/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	type have struct {
		previous model.LoginEvents
		current  *model.LoginEvent
	}
	tts := []struct {
		name string
		have have
		want *Set
	}{
		{
			name: "No previous",
			have: have{
				previous: []*model.LoginEvent{},
				current:  model.MockLoginEvent(model.MockConfig{Suffix: "current"}),
			},
			want: &Set{
				Previous: nil,
				Current:  nil,
				data: []data{
					{
						reason: "NotNoPrevious",
						value:  0,
					},
				},
				result: 0,
			},
		},
		{
			name: "known hash",
			have: have{
				previous: model.LoginEvents{
					model.MockLoginEvent(model.MockConfig{Suffix: "latest", H: 13}),
				},
				current: model.MockLoginEvent(model.MockConfig{Suffix: "current", H: 14}),
			},
			want: &Set{
				Previous: nil,
				Current:  nil,
				data: []data{
					{
						reason: "NotNoPrevious",
						value:  0,
					},
					{
						reason: "KnownHash",
						value:  0,
					},
				},
				result: 0,
			},
		},
		{
			name: "NotKnownDeviceID",
			have: have{
				previous: model.LoginEvents{
					model.MockLoginEvent(model.MockConfig{Suffix: "latest", Hash: "h_other", DeviceID: "d_other", H: 13}),
				},
				current: model.MockLoginEvent(model.MockConfig{Suffix: "current", H: 14}),
			},
			want: &Set{
				Previous: nil,
				Current:  nil,
				data: []data{
					{
						reason: "NotNoPrevious",
						value:  0,
					},
					{
						reason: "NotKnownHash",
						value:  100,
					},
					{
						reason: "NotKnownDeviceID",
						value:  0,
					},
				},
				result: 100,
			},
		},
		{
			name: "knownIP",
			have: have{
				previous: model.LoginEvents{
					model.MockLoginEvent(model.MockConfig{Suffix: "first", H: 12}),
					model.MockLoginEvent(model.MockConfig{Suffix: "second", H: 13}),
					model.MockLoginEvent(model.MockConfig{Suffix: "latest", H: 16, IP: "192.168.1.1"}),
				},
				current: model.MockLoginEvent(model.MockConfig{Suffix: "latest", H: 17, IP: "192.168.1.1"}),
			},
			want: &Set{
				Previous: nil,
				Current:  nil,
				data: []data{
					{
						reason: "NotNoPrevious",
						value:  0,
					},
					{
						reason: "NotKnownHash",
						value:  100,
					},
					{
						reason: "NotKnownDeviceID",
						value:  100,
					},
					{
						reason: "KnownIP",
						value:  0,
					},
				},
				result: 200,
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			c := mockClient(t)
			got, err := c.Run(tt.have.previous, tt.have.current)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
