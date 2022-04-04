package rules

import (
	"errors"
	"geoip/pkg/model"
)

type data struct {
	reason string
	value  int
}

// Set is the rule set
type Set struct {
	Previous model.LoginEvents
	Current  *model.LoginEvent
	data     []data
	result   int
}

// ruleNoPrevious return true if there is no previous loginEvents
func (s *Set) ruleKnownPrevious() bool {
	if s.Previous.IsEmpty() {
		s.data = append(s.data, data{
			reason: "NotKnownPrevious",
			value:  0,
		})
		return false
	}

	s.data = append(s.data, data{
		reason: "KnownPrevious",
		value:  0,
	})
	return true
}

func (s *Set) ruleKnownDeviceID() bool {
	if s.Previous.HasDeviceID(s.Current.DeviceIDHashed) {
		s.data = append(s.data, data{
			reason: "KnownDeviceID",
			value:  0,
		})
		return true
	}

	s.data = append(s.data, data{
		reason: "NotKnownDeviceID",
		value:  100,
	})
	return false
}

// ruleKnownHash return true and value 0 if hash is known, else 100 and false
func (s *Set) ruleKnownHash() bool {
	if s.Previous.HasHash(s.Current.Hash) {
		s.data = append(s.data, data{
			reason: "KnownHash",
			value:  0,
		})
		return true
	}

	s.data = append(s.data, data{
		reason: "NotKnownHash",
		value:  100,
	})
	return false
}

// ruleIsKnownIP return true if ip is present in previous loginEvents.
func (s *Set) ruleKnownIP() {
	if s.Previous.HasIP(s.Current.IP.IPAddr) {
		s.data = append(s.data, data{
			reason: "KnownIP",
			value:  0,
		})
		return
	}
	s.data = append(s.data, data{
		reason: "NotKnownIP",
		value:  50,
	})
}

// ruleKnownCountry add value points country if country is not present in previous loginEvents
func (s *Set) ruleKnownCountry() {
	if s.Previous.HasCountry(s.Current.Location.Country) {
		s.data = append(s.data, data{
			reason: "KnownCountry",
			value:  0,
		})
		return
	}
	s.data = append(s.data, data{
		reason: "NotKnownCountry",
		value:  50,
	})
}

// ruleUserAgent add value points if user-agent browser/os/device is not present in previous loginEvents
func (s *Set) ruleKnownUserAgent() {
	loginEvent := s.Previous.FindMatchingUA(s.Current.UserAgent)
	if loginEvent == nil {
		s.data = append(s.data, data{
			reason: "UserAgent",
			value:  50,
		})
		return
	}

	s.data = append(s.data, data{
		reason: "UserAgent",
		value:  50,
	})
}

func (s *Set) validate() error {
	if s.Current == nil {
		return errors.New("Error: No current loginEvent")
	}

	return nil
}

func (s *Set) resolve() {
	for _, r := range s.data {
		s.result += r.value
	}
	s.Previous = nil
	s.Current = nil
}

// Run runs the rule engine
func (c *Client) Run(previous model.LoginEvents, current *model.LoginEvent) (*Set, error) {
	set := &Set{
		Previous: previous,
		Current:  current,
	}

	if err := set.validate(); err != nil {
		return nil, err
	}

	if !set.ruleKnownPrevious() {
		set.resolve()
		return set, nil
	}

	if set.ruleKnownHash() {
		set.resolve()
		return set, nil
	}

	if set.ruleKnownDeviceID() {
		set.resolve()
		return set, nil
	}

	set.ruleKnownIP()

	set.ruleKnownCountry()

	set.ruleKnownUserAgent()

	set.resolve()

	return set, nil
}
