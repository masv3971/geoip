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

func TestML2Timestamp(t *testing.T) {
	tts := []struct {
		name string
		have int64
		want time.Time
	}{
		{
			name: "15.30",
			have: 55800,
			want: time.Date(0, 0, 0, 15, 30, 00, 00, time.UTC),
		},
		{
			name: "8.30",
			have: 30600,
			want: time.Date(0, 0, 0, 8, 30, 00, 00, time.UTC),
		},
		{
			name: "00.00",
			have: 0,
			want: time.Date(0, 0, 0, 0, 00, 00, 00, time.UTC),
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			got := ml2Timestamp(tt.have)
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
				MockLoginEvent(MockConfig{Suffix: "b", UABrowser: "firefox", UAOS: "linux", UADevice: "laptop"}),
			}.FindMatchingUA(tt.have)

			if tt.want != nil {
				assert.Equal(t, tt.want.UserAgent.Browser.Family, got.UserAgent.Browser.Family)
				assert.Equal(t, tt.want.UserAgent.OS.Family, got.UserAgent.OS.Family)
				assert.Equal(t, tt.want.UserAgent.Device.Family, got.UserAgent.Device.Family)
			} else {
				assert.Equal(t, tt.want, got)
			}
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

func TestIPStats(t *testing.T) {
	tts := []struct {
		name string
		have LoginEvents
		want map[string]int
	}{
		{
			name: "OK",
			have: LoginEvents{
				MockLoginEvent(MockConfig{Suffix: "a", IP: "192.168.0.1"}),
				MockLoginEvent(MockConfig{Suffix: "b", IP: "192.168.0.1"}),
				MockLoginEvent(MockConfig{Suffix: "c", IP: "192.168.0.2"}),
				MockLoginEvent(MockConfig{Suffix: "d", IP: "192.168.0.3"}),
				MockLoginEvent(MockConfig{Suffix: "e", IP: "192.168.0.1"}),
			},
			want: map[string]int{
				"192.168.0.1": 3,
				"192.168.0.2": 1,
				"192.168.0.3": 1,
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.have.IPStats()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNumberOfUniqueIPs(t *testing.T) {
	tts := []struct {
		name string
		have LoginEvents
		want int
	}{
		{
			name: "Two unique ips",
			have: LoginEvents{
				MockLoginEvent(MockConfig{Suffix: "a", IP: "10.0.0.1"}),
				MockLoginEvent(MockConfig{Suffix: "b", IP: "10.0.0.2"}),
			},
			want: 2,
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.have.NumberOfUniqueIPs()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNumberOfUniqueCountries(t *testing.T) {
	tts := []struct {
		name string
		have LoginEvents
		want int
	}{
		{
			name: "Two unique countries",
			have: LoginEvents{
				MockLoginEvent(MockConfig{Suffix: "a", Country: "USA"}),
				MockLoginEvent(MockConfig{Suffix: "b", Country: "Sweden"}),
			},
			want: 2,
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.have.NumberOfUniqueCountries()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHasCountry(t *testing.T) {
	type have struct {
		arg         string
		loginEvents LoginEvents
	}
	tts := []struct {
		name string
		have have
		want bool
	}{
		{
			name: "Using default value",
			have: have{
				loginEvents: LoginEvents{
					MockLoginEvent(MockConfig{Suffix: "a"}),
					MockLoginEvent(MockConfig{Suffix: "a", Country: "USA"}),
				},
				arg: "sweden",
			},
			want: true,
		},
		{
			name: "Using default value",
			have: have{
				loginEvents: LoginEvents{
					MockLoginEvent(MockConfig{Suffix: "a"}),
					MockLoginEvent(MockConfig{Suffix: "a", Country: "usa"}),
				},
				arg: "USA",
			},
			want: true,
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.have.loginEvents.HasCountry(tt.have.arg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNumberOfCountries(t *testing.T) {
	tts := []struct {
		name string
		have LoginEvents
		want int
	}{
		{
			name: "OK",
			have: LoginEvents{
				MockLoginEvent(MockConfig{Suffix: "a", Country: "USA"}),
				MockLoginEvent(MockConfig{Suffix: "b", Country: "England"}),
				MockLoginEvent(MockConfig{Suffix: "c", Country: "ukraine"}),
				MockLoginEvent(MockConfig{Suffix: "d"}),
			},
			want: 4,
		},
	}
	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.have.NumberOfCountries()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHasDeviceID(t *testing.T) {
	type have struct {
		arg         string
		loginEvents LoginEvents
	}
	tts := []struct {
		name string
		have have
		want bool
	}{
		{
			name: "OK",
			have: have{
				loginEvents: LoginEvents{
					MockLoginEvent(MockConfig{Suffix: "a", DeviceID: "test"}),
				},
				arg: "test",
			},
			want: true,
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.have.loginEvents.HasDeviceID(tt.have.arg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHasHash(t *testing.T) {
	tts := []struct {
		name string
		have struct {
			arg string
			le  LoginEvents
		}
		want bool
	}{
		{
			name: "OK",
			have: struct {
				arg string
				le  LoginEvents
			}{
				arg: "test",
				le: LoginEvents{
					MockLoginEvent(MockConfig{Suffix: "a", Hash: "test"}),
					MockLoginEvent(MockConfig{Suffix: "b", Hash: "test"}),
				},
			},

			want: true,
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.have.le.HasHash(tt.have.arg)
			assert.Equal(t, tt.want, got)

		})
	}
}

func TestNumberOfIPs(t *testing.T) {
	tts := []struct {
		name string
		have LoginEvents
		want int
	}{
		{
			name: "OK",
			have: []*LoginEvent{},
			want: 0,
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.have.NumberOfIPs()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParse2ML(t *testing.T) {
	tts := []struct {
		name string
		have *LoginEvent
		want *LoginEvent
	}{
		{
			name: "OK",
			have: MockLoginEvent(MockConfig{Suffix: "a"}),
			want: &LoginEvent{
				ID:          "id_a",
				EppnHashed:  "eppn_a",
				Hash:        "h_abc",
				TimestampML: 46800,
				Timestamp:   time.Date(2022, 2, 23, 13, 0, 0, 0, time.UTC),
				IP: &IP{
					IPAddr:   "10.0.0.1",
					IPAddrML: "167772161",
					ASN: &ASN{
						Number:       1257,
						Organization: "",
					},
					ISP:         &ISP{},
					AnonymousIP: &AnonymousIP{},
				},
				DeviceIDHashed: "d_abc",
				UserAgent: &UserAgent{
					Browser: UserAgentSoftware{
						Family:   "firefox",
						FamilyML: 2,
					},
					OS: UserAgentSoftware{
						Family:   "linux",
						FamilyML: 1,
					},
					Device: UserAgentHardware{
						Family:   "laptop",
						FamilyML: 0,
					},
					Sophisticated: UserAgentSophisticated{},
				},
				Phisheness: &Phisheness{},
				Location: &Location{
					Coordinates: &Coordinates{
						Latitude:  57.648,
						Longitude: 12.5022,
					},
					Country:   "sweden",
					CountryML: 752,
				},
				LoginMethod:   "",
				Fraudulent:    false,
				FraudulentInt: 0,
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.have.Parse2ML()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, tt.have)
		})
	}
}

func TestHasIP(t *testing.T) {
	tts := []struct {
		name string
		have struct {
			arg string
			le  LoginEvents
		}
		want bool
	}{
		{
			name: "OK",
			have: struct {
				arg string
				le  LoginEvents
			}{
				arg: "10.0.0.1",
				le: LoginEvents{
					MockLoginEvent(MockConfig{Suffix: "a"}),
				},
			},
			want: true,
		},
		{
			name: "Not OK",
			have: struct {
				arg string
				le  LoginEvents
			}{
				arg: "10.0.0.2",
				le: LoginEvents{
					MockLoginEvent(MockConfig{Suffix: "a"}),
				},
			},
			want: false,
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.have.le.HasIP(tt.have.arg)
			assert.Equal(t, tt.want, got)
		})
	}
}
