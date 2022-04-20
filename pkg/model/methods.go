package model

import (
	"errors"
	"math/big"
	"strings"
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

func ml2Timestamp(sec int64) time.Time {
	return time.Date(0, 0, 0, 0, 0, int(sec), 0, time.UTC)
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
	switch strings.ToLower(s) {
	case "linux":
		return 1
	case "mac os x", "mac osx":
		return 2
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
	if l == nil {
		return errors.New("Input is nil")
	}
	// ip --> int (as string to prevent overflow)
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

// Equal neglect upper or lower cases in input values, return true -> Sweden, sweden
func Equal(i, j string) bool {
	return strings.ToLower(i) == strings.ToLower(j)
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
	}

	return t
}

// IPStats return a structure of ip's and data about them
func (l LoginEvents) IPStats() map[string]int {
	t := map[string]int{}

	for _, loginEvent := range l {
		if loginEvent == nil {
			break
		}
		if loginEvent.IP.IPAddr == "" {
			break
		}
		t[loginEvent.IP.IPAddr]++
	}
	return t
}

// NumberOfUniqueIPs return int of unique ips
func (l LoginEvents) NumberOfUniqueIPs() int {
	tempMap := map[string]bool{}
	var result int = 0

	for _, loginEvent := range l {
		if _, ok := tempMap[loginEvent.IP.IPAddr]; !ok {
			result++
			tempMap[loginEvent.IP.IPAddr] = true
		}
	}
	return result
}

// NumberOfIPs return all aviable ips in LoginEvents for a eppn
func (l LoginEvents) NumberOfIPs() int {
	var result int = 0

	for _, loginEvent := range l {
		if loginEvent.IP.IPAddr != "" {
			result++
		}
	}
	return result
}

// NumberOfUniqueCountries return int of unique ips
func (l LoginEvents) NumberOfUniqueCountries() int {
	tempMap := map[string]bool{}
	var result int = 0

	for _, loginEvent := range l {
		if _, ok := tempMap[loginEvent.Location.Country]; !ok {
			result++
			tempMap[loginEvent.Location.Country] = true
		}
	}
	return result
}

// NumberOfCountries return all aviable countries in LoginEvents for a eppn
func (l LoginEvents) NumberOfCountries() int {
	var result int = 0

	for _, loginEvent := range l {
		if loginEvent.IP.IPAddr != "" {
			result++
		}
	}
	return result
}

// HasCountry return true if country is found in any of loginEvent
func (l LoginEvents) HasCountry(country string) bool {
	for _, loginEvent := range l {
		if Equal(loginEvent.Location.Country, country) {
			return true
		}
	}
	return false
}

//// HasDeviceID return true if deviceID is found in any of loginEvent
//func (l LoginEvents) HasDeviceID(deviceID string) bool {
//	for _, loginEvent := range l {
//		if Equal(loginEvent.KnownDevice, deviceID) {
//			return true
//		}
//	}
//	return false
//}

// HasHash return true if hash is found in any of loginEvents
func (l LoginEvents) HasHash(hash string) bool {
	for _, loginEvent := range l {
		if Equal(loginEvent.Hash, hash) {
			return true
		}
	}
	return false
}

// FindMatchingUA return loginEvent that match cUA (current loginEvent); --> l[*].UserAgent.* == cUA[*].UserAgent.[*]
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
