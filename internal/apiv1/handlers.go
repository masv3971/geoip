package apiv1

import (
	"context"
	"geoip/internal/statistics"
	"geoip/pkg/model"
	"net"
	"sort"
	"time"
)

// RequestLoginEvent type
type RequestLoginEvent struct {
	Data struct {
		EppnHashed string           `json:"eppn_hashed" validate:"required"`
		ClientIP   net.IP           `json:"client_ip" validate:"required"`
		DeviceID   bool           `json:"device_id" validate:"required"`
		UserAgent  *model.UserAgent `json:"user_agent"`
	} `json:"data" validate:"required"`
}

// ReplyLoginEvent reply type for LoginEvent
type ReplyLoginEvent struct {
	Status bool `json:"status"`
}

// HandlerLoginEvent return a loginEvent, or error
func (c *Client) HandlerLoginEvent(ctx context.Context, indata *RequestLoginEvent) (*ReplyLoginEvent, error) {
	var (
		enterpriseISP         *model.ISP
		enterpriseAnonymousIP *model.AnonymousIP
	)

	city, err := c.maxmind.City(indata.Data.ClientIP)
	if err != nil {
		return nil, err
	}
	asn, err := c.maxmind.ASN(indata.Data.ClientIP)
	if err != nil {
		return nil, err
	}

	if c.config.MaxMind.Enterprise {
		isp, err := c.maxmind.ISP(indata.Data.ClientIP)
		if err != nil {
			return nil, err
		}
		enterpriseISP = &model.ISP{
			ASN:            isp.AutonomousSystemNumber,
			ASOrganization: isp.AutonomousSystemOrganization,
			ISP:            isp.ISP,
			Organization:   isp.Organization,
		}
		anonymousIP, err := c.maxmind.AnonymousIP(indata.Data.ClientIP)
		if err != nil {
			return nil, err
		}

		enterpriseAnonymousIP = &model.AnonymousIP{
			IsAnonymous:        anonymousIP.IsAnonymous,
			IsAnonymousVPN:     anonymousIP.IsAnonymousVPN,
			IsHostingProvider:  anonymousIP.IsHostingProvider,
			IsPublicProxy:      anonymousIP.IsPublicProxy,
			IsResidentialProxy: anonymousIP.IsResidentialProxy,
			IsTorExitNode:      anonymousIP.IsTorExitNode,
		}
	}

	loginEventCurrent := &model.LoginEvent{
		EppnHashed:     indata.Data.EppnHashed,
		KnownDevice: indata.Data.DeviceID,
		Timestamp:      time.Now(),
		Location: &model.Location{
			Coordinates: &model.Coordinates{Latitude: city.Location.Latitude, Longitude: city.Location.Longitude},
			City:        city.City.Names["en"],
			Country:     city.Country.Names["en"],
			Continent:   city.Continent.Names["en"],
		},
		IP: &model.IP{
			IPAddr: indata.Data.ClientIP.String(),
			ASN: &model.ASN{
				Number:       asn.AutonomousSystemNumber,
				Organization: asn.AutonomousSystemOrganization,
			},
			ISP:         enterpriseISP,
			AnonymousIP: enterpriseAnonymousIP,
		},
		UserAgent: indata.Data.UserAgent,
	}

	if err := loginEventCurrent.Parse2ML(); err != nil {
		return nil, err
	}

	if err := loginEventCurrent.AddHash(); err != nil {
		return nil, err
	}

	if _, err := c.store.AddLoginEvent(ctx, loginEventCurrent); err != nil {
		return nil, err
	}

	return &ReplyLoginEvent{Status: true}, nil
}

// RequestStatsOverview reply type
type RequestStatsOverview struct{}

// ReplyStatsOverview reply type for HandlerStatsCollection
type ReplyStatsOverview struct {
	model.StatsOverviewDocs `json:"overview"`
}

// HandlerStatsOverview handler for stats collection
func (c *Client) HandlerStatsOverview(ctx context.Context, indata *RequestStatsOverview) (*ReplyStatsOverview, error) {
	all, err := c.store.GetLoginEventsAll(ctx)
	if err != nil {
		return nil, err
	}

	statsDocuments := model.StatsOverviewDocs{}

	for eppn, loginEvents := range all {
		statsDocuments = append(statsDocuments, model.StatsOverviewDoc{
			EPPNHashed:              eppn,
			NumbnerOfLoginEvents:    len(loginEvents),
			NumberOfCountries:       loginEvents.NumberOfCountries(),
			NumberOfUniqueCountries: loginEvents.NumberOfUniqueCountries(),
			NumberOfIPs:             loginEvents.NumberOfIPs(),
			NumberOfUniqueIPs:       loginEvents.NumberOfUniqueIPs(),
		})
	}

	sort.Sort(model.StatsOverviewDocsSortByOccurrences(statsDocuments))

	return &ReplyStatsOverview{
		StatsOverviewDocs: statsDocuments,
	}, nil
}

// RequestStatsEppnLong input type
type RequestStatsEppnLong struct {
	EppnHashed string `uri:"eppn"`
}

// ReplyStatsEppnLong output type
type ReplyStatsEppnLong struct {
	model.LoginEvents
}

// HandlerStatsEppnLong return loginEvent for a Eppn
func (c *Client) HandlerStatsEppnLong(ctx context.Context, indata *RequestStatsEppnLong) (*ReplyStatsEppnLong, error) {
	loginEvents, err := c.store.GetLoginEvents(ctx, indata.EppnHashed)
	if err != nil {
		return nil, err
	}

	return &ReplyStatsEppnLong{loginEvents}, nil
}

// RequestStatsEppnSpecific input type
type RequestStatsEppnSpecific struct {
	EppnHashed string `uri:"eppn"`
}

// ReplyStatsEppnSpecific output type
type ReplyStatsEppnSpecific struct {
	statistics.StatsData
}

// HandlerStatsEppnSpecific return summery for one eppn
func (c *Client) HandlerStatsEppnSpecific(ctx context.Context, indata *RequestStatsEppnSpecific) (*ReplyStatsEppnSpecific, error) {
	loginEvents, err := c.store.GetLoginEvents(ctx, indata.EppnHashed)
	if err != nil {
		return nil, err
	}

	sData, err := statistics.NewSpecific(loginEvents)
	if err != nil {
		return nil, err
	}

	return &ReplyStatsEppnSpecific{sData}, nil
}
