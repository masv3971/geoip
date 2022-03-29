package model

import (
	"fmt"
	"math/big"
	"time"

	"github.com/biter777/countries"
	"inet.af/netaddr"
)

func (f Fraudulent) String() string {
	if f == true {
		return "true"
	}
	return "false"
}

func ip2ML(ip string) (string, error) {
	IP := netaddr.MustParseIP(ip)
	ipBinary, err := IP.MarshalBinary()
	if err != nil {
		return "", err
	}

	ipInt := big.NewInt(0)
	ipInt.SetBytes(ipBinary)
	return ipInt.String(), nil
}

func timestamp2ML(ts time.Time) int64 {
	h, m, s := ts.Clock()
	hSeconds := int64(h) * 60 * 60
	mSeconds := int64(m) * 60

	return hSeconds + mSeconds + int64(s)
}

func country2ML(countryName string) int {
	return int(countries.ByName(countryName))
}

func uaBrowserFamily2ML(s string) int {
	switch s {
	case "chrome", "Chrome":
		return 1
	case "firefox", "Firefox":
		return 2
	default:
		return 0
	}
}

func uaOSFamily2ML(s string) int {
	switch s {
	case "Mac OS X", "mac os x", "mac os X", "mac osx", "mac osX":
		return 1
	default:
		return 0
	}
}

func usDeviceFamily2ML(s string) int {
	switch s {
	case "":
		return 1
	default:
		return 0
	}
}

// Parse2ML parse a loginEvent fields used by marchine learning models
func (l *LoginEvent) Parse2ML() error {
	// ip --> int (as string to prevent )
	ipInt, err := ip2ML(l.IP.IPAddr)
	if err != nil {
		return err
	}
	l.IP.IPAddrML = ipInt

	// timestamp --> int
	l.TimestampML = timestamp2ML(l.Timestamp)

	// country --> int
	l.Location.CountryML = country2ML(l.Location.Country)

	// userAgent browser family --> int
	l.UserAgent.Browser.FamilyML = uaBrowserFamily2ML(l.UserAgent.Browser.Family)

	// userAgent OS family --> int
	l.UserAgent.OS.FamilyML = uaOSFamily2ML(l.UserAgent.OS.Family)

	// userAgent Device family --> int
	l.UserAgent.Device.FamilyML = usDeviceFamily2ML(l.UserAgent.Device.Family)

	return nil
}

// IsEmpty return true if LoginEvents is of length 0, else false
func (l LoginEvents) IsEmpty() bool {
	return len(l) == 0
}

// HasIP return true if ip is found in any of loginEvent
func (l LoginEvents) HasIP(ip string) bool {
	for _, loginEvent := range l {
		if loginEvent.IP.IPAddr == ip {
			return true
		}
	}
	return false
}

// CountriesStat return a compilated structure of countries and data about them
func (l LoginEvents) CountriesStat() map[string]int {
	t := map[string]int{}
	for _, loginEvent := range l {
		if loginEvent == nil {
			break
		}
		if loginEvent.Location.Country == "" {
			break
		}
		t[loginEvent.Location.Country]++
		fmt.Println("location country", t[loginEvent.Location.Country], loginEvent.Location.Country)
	}

	return t
}

// HasCountry return true if country is found in any of loginEvent
func (l LoginEvents) HasCountry(country string) bool {
	for _, loginEvent := range l {
		if loginEvent.Location.Country == country {
			return true
		}
	}
	return false
}

// HasDeviceID return true if deviceID is found in any of loginEvent
func (l LoginEvents) HasDeviceID(deviceID string) bool {
	for _, loginEvent := range l {
		if loginEvent.DeviceIDHashed == deviceID {
			return true
		}
	}
	return false
}

// HasHash return true if hash is found in any of loginEvents
func (l LoginEvents) HasHash(hash string) bool {
	for _, loginEvent := range l {
		if loginEvent.Hash == hash {
			return true
		}
	}
	return false
}

// FindMatchingUA return loginEvent that match cUA (current loginEvent)
func (l LoginEvents) FindMatchingUA(cUA *UserAgent) *LoginEvent {
	for _, loginEvent := range l {
		if loginEvent.UserAgent.Browser.Family == cUA.Browser.Family {
			if loginEvent.UserAgent.Device.Family == cUA.Device.Family {
				if loginEvent.UserAgent.OS.Family == cUA.OS.Family {
					return loginEvent
				}
			}
		}
	}
	return nil
}
