package model

import (
	"fmt"
	"strings"
	"time"
)

func mockLocation(country string) *Coordinates {
	switch strings.ToLower(country) {
	case "sweden":
		return &Coordinates{
			Latitude:  57.648,
			Longitude: 12.5022,
		}
	case "usa":
		return &Coordinates{
			Latitude:  44.500000,
			Longitude: -89.500000,
		}
	default:
		return nil
	}
}

// MockConfig configs MockLoginEvent
type MockConfig struct {
	Suffix, Hash, IP, Country, UABrowser, UAOS, UADevice string
	H, M, S                                              int
	ASN                                                  uint
	KnownDevice                                          bool
}

func (c *MockConfig) defaultConfig() {
	if c.Hash == "" {
		c.Hash = "h_abc"
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

func (c *MockConfig) toLower() {
	c.Hash = strings.ToLower(c.Hash)
	c.IP = strings.ToLower(c.IP)
	c.UABrowser = strings.ToLower(c.UABrowser)
	c.UADevice = strings.ToLower(c.UADevice)
	c.UAOS = strings.ToLower(c.UAOS)
}

// MockLoginEvent mocks loginEvent with a sane default
func MockLoginEvent(c MockConfig) *LoginEvent {
	c.defaultConfig()
	c.toLower()

	le := &LoginEvent{
		ID:          fmt.Sprintf("id_%s", c.Suffix),
		EppnHashed:  fmt.Sprintf("eppn_%s", c.Suffix),
		Hash:        c.Hash,
		KnownDevice: c.KnownDevice,
		Location: &Location{
			Country:     c.Country,
			Coordinates: mockLocation(c.Country),
		},
		Timestamp: time.Date(2022, 2, 23, c.H, c.M, c.S, 0, time.UTC),
		IP: &IP{
			IPAddr:      c.IP,
			IPAddrML:    "",
			ASN:         &ASN{Number: c.ASN, Organization: ""},
			ISP:         &ISP{},
			AnonymousIP: &AnonymousIP{},
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
		Phisheness:    &Phisheness{},
	}

	return le
}
