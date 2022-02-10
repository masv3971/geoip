package model

import (
	"strconv"
	"time"

	"github.com/mitchellh/hashstructure/v2"
)

// LoginEvent keep record of each login event
type LoginEvent struct {
	ID        string     `json:"uid,omitempty" bson:"_id,omitempty" hash:"-"`
	Eppn      string     `json:"eppn" bson:"eppn"`
	Hash      string     `json:"hash" bson:"hash" hash:"-"`
	TimeStamp time.Time  `json:"time_stamp,omitempty" hash:"-"`
	IP        *IP        `json:"ip,omitempty"`
	Travel    *Travel    `json:"travel,omitempty"`
	DeviceID  string     `json:"device_id" bson:"device_id"`
	UserAgent *UserAgent `json:"user_agent" bson:"user_agent"`
	Phishing  Phishing   `json:"phishing" bson:"phishing" hash:"-"`
}

// IP keeps track of ip
type IP struct {
	IPAddr      string       `json:"ip_addr,omitempty" bson:"ip_addr"`
	Country     string       `json:"country,omitempty" bson:"country"`
	City        string       `json:"city,omitempty" bson:"city"`
	Coordinates *Coordinates `json:"location,omitempty" bson:"coordinates"`
	ASN         *ASN         `json:"asn" bson:"asn"`
	ISP         *ISP         `json:"isp" bson:"isp"`
	AnonymousIP *AnonymousIP `json:"anonymous_ip" bson:"anonymous_ip"`
}

// ASN Autonomous System Number, info like asn, organization etc
type ASN struct {
	Number       uint   `json:"number"`
	Organization string `json:"organization"`
}

// ISP internet service provider infomation
type ISP struct {
	ASN            uint   `json:"asn" bson:"asn"`
	ASOrganization string `json:"as_organization" bson:"as_organization"`
	ISP            string `json:"isp" bson:"isp"`
	Organization   string `json:"organization" bson:"organization"`
}

// AnonymousIP information about if the ip has any clocking devices
type AnonymousIP struct {
	IsAnonymous        bool `json:"is_anonymous" bson:"is_anonymous"`
	IsAnonymousVPN     bool `json:"is_anonymous_vpn" bson:"is_anonymous_vpn"`
	IsHostingProvider  bool `json:"is_hosting_provider" bson:"is_hosting_provider"`
	IsPublicProxy      bool `json:"is_public_proxy" bson:"is_public_proxy"`
	IsResidentialProxy bool `json:"is_residential_proxy" bson:"is_residential_proxy"`
	IsTorExitNode      bool `json:"is_tor_exit_node" bson:"is_tor_exit_node"`
}

// Coordinates keeps track of each IPs location
type Coordinates struct {
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

// Travel contains data for a travel
type Travel struct {
	Distance           float64       `json:"Distance" bson:"distance" hash:"-"`
	DistanceUnit       string        `json:"unit" bson:"distance_unit"`
	TravelDuration     time.Duration `json:"travel_duration" bson:"time_of_travel" hash:"-"`
	TravelDurationUnit string        `json:"travel_duration_unit" bson:"travel_duration_unit"`
	TravelSpeed        float64       `json:"travel_speed" bson:"travel_speed" hash:"-"`
	TravelSpeedUnit    string        `json:"travel_speed_unit" bson:"travel_speed_unit"`
	IsTravelImpossible bool          `json:"is_travel_impossible" bson:"is_travel_impossible" hash:"-"`
}

// UserAgent contains data from users user-agent
type UserAgent struct {
	Browser struct {
		Family        string `json:"family" bson:"family"`
		Version       []int  `json:"version" bson:"version"`
		VersionString string `json:"version_string" bson:"version_string"`
	} `json:"browser"`
	OS struct {
		Family        string `json:"family" bson:"family"`
		Version       []int  `json:"version" bson:"version"`
		VersionString string `json:"version_string" bson:"version_string"`
	} `json:"os"`
	Device struct {
		Family string `json:"family" bson:"family"`
		Brand  string `json:"brand" bson:"brand"`
		Model  string `json:"model" bson:"model"`
	} `json:"device"`
	Sophisticated struct {
		IsMobile       bool `json:"is_mobile" bson:"is_mobile"`
		IsTablet       bool `json:"is_tablet" bson:"is_tablet"`
		IsPC           bool `json:"is_pc" bson:"is_pc"`
		IsTouchCapable bool `json:"is_touch_capable" bson:"is_touch_capable"`
		IsBot          bool `json:"is_bot" bson:"is_bot"`
	} `json:"sophisticated" bson:"sophisticated"`
}

// Phishing contains data about phishing
type Phishing struct {
	Score  int    `json:"score" bson:"score"`
	Reason string `json:"reason" bson:"reason"`
}

// AddHash make a hash of LoginEvent
func (l *LoginEvent) AddHash() error {
	hash, err := hashstructure.Hash(*l, hashstructure.FormatV2, nil)
	if err != nil {
		return err
	}

	l.Hash = strconv.FormatUint(hash, 10)
	return nil
}
