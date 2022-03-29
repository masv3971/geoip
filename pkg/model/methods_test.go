package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetLatest(t *testing.T) {
	got := LoginEvents{
		MockLoginEvent(MockConfig{Suffix: "a", H: 13}),
		MockLoginEvent(MockConfig{Suffix: "b", H: 8}),
		MockLoginEvent(MockConfig{Suffix: "latest", H: 15}),
	}.GetLatest()
	assert.Equal(t, MockLoginEvent(MockConfig{Suffix: "latest", H: 15}), got)

}

func TestFraudulentString(t *testing.T) {
	tts := []struct {
		name string
		have Fraudulent
		want string
	}{
		{
			name: "OK-false",
			have: false,
			want: "false",
		},
		{
			name: "OK-true",
			have: true,
			want: "true",
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.have.String()
			assert.Equal(t, tt.want, got)

		})
	}
}

func TestIP2ML(t *testing.T) {
	tts := []struct {
		name string
		have string
		want string
	}{
		{
			name: "ipv4",
			have: "192.168.0.1",
			want: "3232235521",
		},
		{
			name: "ipv4",
			have: "192.168.0.11",
			want: "3232235531",
		},
		{
			name: "ipv6",
			have: "2a00:801:23b:d280::1",
			want: "55827738162162749090732223022403420161",
		},
		{
			name: "ipv6",
			have: "2a00:801:23b:d280::a",
			want: "55827738162162749090732223022403420170",
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ip2ML(tt.have)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTimestamp2ML(t *testing.T) {
	tts := []struct {
		name string
		have time.Time
		want int64
	}{
		{
			name: "15.30",
			have: time.Date(2022, 2, 2, 15, 30, 00, 00, time.UTC),
			want: 55800,
		},
		{
			name: "8.30",
			have: time.Date(2022, 2, 2, 8, 30, 00, 00, time.UTC),
			want: 30600,
		},
		{
			name: "00.00",
			have: time.Date(2022, 2, 2, 0, 00, 00, 00, time.UTC),
			want: 0,
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			got := timestamp2ML(tt.have)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCountry2ML(t *testing.T) {
	tts := []struct {
		name string
		have string
		want int
	}{
		{
			name: "Sweden",
			have: "Sweden",
			want: 752,
		},
		{
			name: "sweden",
			have: "sweden",
			want: 752,
		},
		{
			name: "russia",
			have: "russia",
			want: 643,
		},
		{
			name: "england",
			have: "england",
			want: 826,
		},
		{
			name: "norway",
			have: "norway",
			want: 578,
		},
		{
			name: "finland",
			have: "finland",
			want: 246,
		},
		{
			name: "denmark",
			have: "denmark",
			want: 208,
		},
		{
			name: "germany",
			have: "germany",
			want: 276,
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			got := country2ML(tt.have)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUABrowserFamily2ML(t *testing.T) {
	tts := []struct {
		name string
		have string
		want int
	}{
		{
			name: "chrome",
			have: "chrome",
			want: 1,
		},
		{
			name: "firefox",
			have: "firefox",
			want: 2,
		},
		{
			name: "unknown",
			have: "unknown-browser",
			want: 0,
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			got := uaBrowserFamily2ML(tt.have)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUAOSFamily2ML(t *testing.T) {
	tts := []struct {
		name string
		have string
		want int
	}{
		{
			name: "mac",
			have: "Mac os x",
			want: 1,
		},
		{
			name: "unknown",
			have: "unknown-browser",
			want: 0,
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			got := uaOSFamily2ML(tt.have)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFindMatchingUA(t *testing.T) {
	tts := []struct {
		name string
		have *UserAgent
		want *LoginEvent
	}{
		{
			name: "true",
			have: &UserAgent{
				Browser: UserAgentSoftware{Family: "firefox"},
				OS:      UserAgentSoftware{Family: "linux"},
				Device:  UserAgentHardware{Family: "laptop"},
			},
			want: MockLoginEvent(MockConfig{Suffix: "latest", UABrowser: "firefox", UAOS: "linux", UADevice: "laptop"}),
		},
		{
			name: "false",
			have: &UserAgent{
				Browser: UserAgentSoftware{Family: "firefox"},
				OS:      UserAgentSoftware{Family: "windows"},
				Device:  UserAgentHardware{Family: "laptop"},
			},
			want: nil,
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			got := LoginEvents{
				MockLoginEvent(MockConfig{Suffix: "a", UABrowser: "explorer", UAOS: "windows", UADevice: "laptop"}),
				MockLoginEvent(MockConfig{Suffix: "b", UABrowser: "firefox", UAOS: "windows", UADevice: "laptop"}),
			}.FindMatchingUA(tt.have)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCountriesStat(t *testing.T) {
	tts := []struct {
		name string
		have LoginEvents
		want map[string]int
	}{
		{
			name: "OK",
			have: LoginEvents{
				MockLoginEvent(MockConfig{Suffix: "b", Country: "Sweden"}),
				MockLoginEvent(MockConfig{Suffix: "c", Country: "Sweden"}),
				MockLoginEvent(MockConfig{Suffix: "d", Country: "Finland"}),
				MockLoginEvent(MockConfig{Suffix: "a", Country: "Sweden"}),
				MockLoginEvent(MockConfig{Suffix: "a", Country: "Denmark"}),
			},
			want: map[string]int{
				"Sweden":  3,
				"Finland": 1,
				"Denmark": 1,
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.have.CountriesStat()
			assert.Equal(t, tt.want, got)
		})
	}
}
