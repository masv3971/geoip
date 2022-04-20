package statistics

import (
	"geoip/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockLoginEventML(t *testing.T, conf model.MockConfig) *model.LoginEvent {
	le := model.MockLoginEvent(conf)
	err := le.Parse2ML()
	assert.NoError(t, err)

	return le
}
func TestNewSpecific(t *testing.T) {
	tts := []struct {
		name string
		have model.LoginEvents
		want StatsData
	}{
		{
			name: "",
			have: model.LoginEvents{
				mockLoginEventML(t, model.MockConfig{Suffix: "a", IP: "10.0.0.2", UABrowser: "firefox"}),
				mockLoginEventML(t, model.MockConfig{Suffix: "b", IP: "10.0.0.3", UABrowser: "safari"}),
				mockLoginEventML(t, model.MockConfig{Suffix: "c", IP: "10.0.0.4", UABrowser: "chrome"}),
				mockLoginEventML(t, model.MockConfig{Suffix: "d", IP: "10.0.0.4", UABrowser: "mozilla"}),
			},
			want: map[string]*model.StatsData{
				"ip": {
					Len:               4,
					Entropy:           0,
					StandardDeviation: 0.82915619758885,
				},
				"user_agent_device": {
					Len:               4,
					Entropy:           0,
					StandardDeviation: 0,
				},
				"user_agent_os": {
					Len:               4,
					Entropy:           0,
					StandardDeviation: 0,
				},
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSpecific(tt.have)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
