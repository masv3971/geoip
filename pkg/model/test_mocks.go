package model

import (
	"fmt"
	"time"
)

func mockLocation(country string) *Coordinates {
	switch country {
	case "sweden":
		return &Coordinates{
			Latitude:  57.648,
			Longitude: 12.5022,
		}
	case "usa":
		return &Coordinates{
			Latitude:  59.3274,
			Longitude: 18.0653,
		}
	default:
		return nil
	}
}

// MockConfig configs MockLoginEvent
type MockConfig struct {
	Suffix, Hash, DeviceID, IP, Country, UABrowser, UAOS, UADevice string
	H, M, S                                                        int
	ASN                                                            uint
}

func (c MockConfig) defaultConfig() {
	if c.Hash == "" {
		c.Hash = "h_abc"
	}
	if c.DeviceID == "" {
		c.DeviceID = "d_abc"
	}
	if c.IP == "" {
		c.IP = "10.0.0.1"
	}
	if c.Country == "" {
		c.Country = "sweden"
	}
	if c.UABrowser == "" {
		c.UABrowser = "firefox"
	}
	if c.UAOS == "" {
		c.UAOS = "linux"
	}
	if c.UADevice == "" {
		c.UADevice = "laptop"
	}
	if c.H == 0 {
		c.H = 13
	}
	if c.ASN == 0 {
		c.ASN = 1257
	}
}

// MockLoginEvent mocks loginEvent with a sane default
func MockLoginEvent(c MockConfig) *LoginEvent {
	c.defaultConfig()

	le := &LoginEvent{
		ID:       fmt.Sprintf("id_%s", c.Suffix),
		EppnHashed:     fmt.Sprintf("eppn_%s", c.Suffix),
		Hash:     c.Hash,
		DeviceIDHashed: c.DeviceID,
		Location: &Location{
			Country:     c.Country,
			Coordinates: mockLocation(c.Country),
		},
		Timestamp: time.Date(2022, 2, 23, c.H, c.M, c.S, 0, time.UTC),
		IP: &IP{
			IPAddr: c.IP,
			ASN: &ASN{
				Number:       c.ASN,
				Organization: "",
			},
		},
		UserAgent: &UserAgent{
			Browser:       UserAgentSoftware{Family: c.UABrowser},
			OS:            UserAgentSoftware{Family: c.UAOS},
			Device:        UserAgentHardware{Family: c.UADevice},
			Sophisticated: UserAgentSophisticated{},
		},
		LoginMethod:   "",
		Fraudulent:    false,
		FraudulentInt: 0,
	}

	return le
}
