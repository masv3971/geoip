package model

import (
	"sort"
	"strconv"
	"time"

	"github.com/mitchellh/hashstructure/v2"
)

// LoginEvents array of *LoginEvent
type LoginEvents []*LoginEvent

// Fraudulent type bool
type Fraudulent bool

// LoginEvent keep record of each login event
type LoginEvent struct {
	ID             string      `json:"uid,omitempty" bson:"_id,omitempty" hash:"-"`
	EppnHashed     string      `json:"eppn_hashed" bson:"eppn_hashed"`
	Hash           string      `json:"hash" bson:"hash" hash:"-"`
	Timestamp      time.Time   `json:"timestamp,omitempty" hash:"-"`
	TimestampML    int64       `json:"timestamp_ml" bson:"timestamp_ml" hash:"-"`
	IP             *IP         `json:"ip,omitempty"`
	DeviceIDHashed string      `json:"device_id_hashed" bson:"device_id_hashed"`
	UserAgent      *UserAgent  `json:"user_agent" bson:"user_agent"`
	Phisheness     *Phisheness `json:"phishing" bson:"phishing" hash:"-"`
	Location       *Location   `json:"location" bson:"location"`
	LoginMethod    string      `json:"login_method" bson:"login_method"`
	Fraudulent     Fraudulent  `json:"fraudulent" bson:"fraudulent"`
	FraudulentInt  int         `json:"fraudulent_ml" bson:"fraudulent_ml"`
}

// IP keeps track of ip
type IP struct {
	IPAddr      string       `json:"ip_addr,omitempty" bson:"ip_addr"`
	IPAddrML    string       `json:"ip_addr_ML" bson:"ip_addr_ml"`
	ASN         *ASN         `json:"asn" bson:"asn"`
	ISP         *ISP         `json:"isp" bson:"isp"`
	AnonymousIP *AnonymousIP `json:"anonymous_ip" bson:"anonymous_ip"`
}

// Location keep location data
type Location struct {
	Coordinates  *Coordinates `json:"coordinates" bson:"coordinates"`
	City         string       `json:"city,omitempty" bson:"city"`
	CityML       int          `json:"city_ml" bson:"city_ml"`
	Country      string       `json:"country,omitempty" bson:"country"`
	CountryML    int          `json:"country_ml" bson:"country_ml"`
	Continent    string       `json:"continent" bson:"continent"`
	ContinentsML int          `json:"continents_ml" bson:"continents_ml"`
}

// Coordinates keeps track of each IPs location
type Coordinates struct {
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

// ASN Autonomous System Number
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
	IsAnonymous          bool `json:"is_anonymous" bson:"is_anonymous"`
	IsAnonymousML        int  `json:"is_anonymous_ml" bson:"is_anonymous_ml"`
	IsAnonymousVPN       bool `json:"is_anonymous_vpn" bson:"is_anonymous_vpn"`
	IsAnonymousVPNML     int  `json:"is_anonymous_vpn_ml" bson:"is_anonymous_vpn_ml"`
	IsHostingProvider    bool `json:"is_hosting_provider" bson:"is_hosting_provider"`
	IsHostingProviderML  int  `json:"is_hosting_provider_ml" bson:"is_hosting_provider_ml"`
	IsPublicProxy        bool `json:"is_public_proxy" bson:"is_public_proxy"`
	IsPublicProxyML      int  `json:"is_public_proxy_ml" bson:"is_public_proxy_ml"`
	IsResidentialProxy   bool `json:"is_residential_proxy" bson:"is_residential_proxy"`
	IsResidentialProxyML int  `json:"is_residential_proxy_ml" bson:"is_residential_proxy_ml"`
	IsTorExitNode        bool `json:"is_tor_exit_node" bson:"is_tor_exit_node"`
	IsTorExitNodeML      int  `json:"is_tor_exit_node_ml" bson:"is_tor_exit_node_ml"`
}

// UserAgentSoftware represents user-agent software
type UserAgentSoftware struct {
	Family        string `json:"family" bson:"family"`
	FamilyML      int    `json:"family_ml" bson:"family_ml"`
	Version       []int  `json:"version" bson:"version"`
	VersionString string `json:"version_string" bson:"version_string"`
}

// UserAgentHardware represents user-agent hardware
type UserAgentHardware struct {
	Family   string `json:"family" bson:"family"`
	FamilyML int    `json:"family_ml" bson:"family_ml"`
	Brand    string `json:"brand" bson:"brand"`
	Model    string `json:"model" bson:"model"`
}

// UserAgentSophisticated represents user-agent sophisticated values
type UserAgentSophisticated struct {
	IsMobile         bool `json:"is_mobile" bson:"is_mobile"`
	IsMobileML       int  `json:"is_mobile_ml" bson:"is_mobile_ml"`
	IsTablet         bool `json:"is_tablet" bson:"is_tablet"`
	IsTabletML       int  `json:"is_tablet_ml" bson:"is_tablet_ml"`
	IsPC             bool `json:"is_pc" bson:"is_pc"`
	IsPCML           int  `json:"is_pc_ml" bson:"is_pc_ml"`
	IsTouchCapable   bool `json:"is_touch_capable" bson:"is_touch_capable"`
	IsTouchCapableML int  `json:"is_touch_capable_ml" bson:"is_touch_capable_ml"`
	IsBot            bool `json:"is_bot" bson:"is_bot"`
	IsBotML          int  `json:"is_bot_ml" bson:"is_bot_ml"`
}

// UserAgent contains data from users user-agent
type UserAgent struct {
	Browser       UserAgentSoftware      `json:"browser"`
	OS            UserAgentSoftware      `json:"os"`
	Device        UserAgentHardware      `json:"device"`
	Sophisticated UserAgentSophisticated `json:"sophisticated" bson:"sophisticated"`
}

// Phisheness contains data about phishing
type Phisheness struct {
	Score        int    `json:"score" bson:"score"`
	Reason       string `json:"reason" bson:"reason"`
	LoginEventID string `json:"login_event_id" bson:"login_event_id"`
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

// GetLatest return the most recent loginEvent from LoginEvents
func (loginEvents LoginEvents) GetLatest() *LoginEvent {
	sort.Sort(LoginEventsSortByTimestamp(loginEvents))
	return loginEvents[0]
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

// StatsOverviewDoc holds the document for stats endpoint
type StatsOverviewDoc struct {
	EPPNHashed              string `json:"eppn_hashed"`
	NumbnerOfLoginEvents    int    `json:"number_of_login_events"`
	NumberOfCountries       int    `json:"number_of_countries"`
	NumberOfUniqueCountries int    `json:"number_of_unique_countries"`
	NumberOfIPs             int    `json:"number_of_ips"`
	NumberOfUniqueIPs       int    `json:"number_of_unique_ips"`
}

type StatsOverviewDocs []StatsOverviewDoc

// StatsData hold statsData
type StatsData struct {
	Len               int     `json:"number_of_elements"`
	Entropy           float64 `json:"entropy"`
	StandardDeviation float64 `json:"standardDeviation"`
}

// StatsSpecificDoc holds data for specific datapoints for a eppn
type StatsSpecificDoc struct {
	IP        StatsData `json:"ip"`
	UserAgent StatsData `json:"user_agent"`
	Country   StatsData `json:"country"`
}
