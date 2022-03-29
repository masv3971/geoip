package apiv1

import (
	"context"
	"geoip/pkg/model"
	"net"
	"time"
)

// RequestLoginEvent type
type RequestLoginEvent struct {
	Data struct {
		EppnHashed string           `json:"eppn_hashed" validate:"required"`
		ClientIP   net.IP           `json:"client_ip" validate:"required"`
		DeviceID   string           `json:"device_id" validate:"required"`
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

	//loginEventsPrevious, err := c.store.GetLoginEvents(ctx, indata.Data.Eppn)
	//if err != nil {
	//	if err != mongo.ErrNoDocuments {
	//		return nil, err
	//	}
	//}

	loginEventCurrent := &model.LoginEvent{
		EppnHashed:     indata.Data.EppnHashed,
		DeviceIDHashed: indata.Data.DeviceID,
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

	//tr, err := tribunal.New(ctx)
	//if err != nil {
	//	return nil, err
	//}

	//phisheness, err := tr.Resolve(ctx, loginEventCurrent, loginEventsPrevious)
	//if err != nil {
	//	return nil, err
	//}

	//loginEventCurrent.Phisheness = phisheness

	if _, err := c.store.AddLoginEvent(ctx, loginEventCurrent); err != nil {
		return nil, err
	}

	return &ReplyLoginEvent{Status: true}, nil
}

// RequestStatsOverview reply type
type RequestStatsOverview struct{}

// ReplyStatsOverview reply type for HandlerStatsCollection
type ReplyStatsOverview struct {
	model.StatsDocuments `json:"stats_documents"`
}

// HandlerStatsOverview handler for stats collection
func (c *Client) HandlerStatsOverview(ctx context.Context, indata *RequestStatsOverview) (*ReplyStatsOverview, error) {
	all, err := c.store.GetLoginEventsAll(ctx)
	if err != nil {
		return nil, err
	}

	statsDocuments := model.StatsDocuments{}

	for eppn, v := range all {
		statsDocuments = append(statsDocuments, model.StatsDocument{
			EPPNHashed:           eppn,
			NumbnerOfLoginEvents: len(v),
			Countries:            v.CountriesStat(),
		})
	}

	return &ReplyStatsOverview{
		StatsDocuments: statsDocuments,
	}, nil
}
