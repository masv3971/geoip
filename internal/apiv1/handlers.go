package apiv1

import (
	"context"
	"geoip/pkg/model"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// RequestLoginEvent type
type RequestLoginEvent struct {
	Data struct {
		Eppn      string           `json:"eppn" validate:"required"`
		ClientIP  net.IP           `json:"client_ip" validate:"required"`
		DeviceID  string           `json:"device_id" validate:"required"`
		UserAgent *model.UserAgent `json:"user_agent"`
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

	loginEventLatest, err := c.store.GetLatestLoginEvent(ctx, indata.Data.Eppn)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, err
		}
	}

	travel, err := c.traveler.Travel(ctx, city.Location.Latitude, city.Location.Longitude, loginEventLatest)
	if err != nil {
		return nil, err
	}

	loginEventCurrent := &model.LoginEvent{
		Eppn:      indata.Data.Eppn,
		TimeStamp: time.Now(),
		IP: &model.IP{
			IPAddr:  indata.Data.ClientIP.String(),
			Country: city.Country.Names["en"],
			City:    city.City.Names["en"],
			Coordinates: &model.Coordinates{
				Latitude:  city.Location.Latitude,
				Longitude: city.Location.Longitude,
			},
			ASN: &model.ASN{
				Number:       asn.AutonomousSystemNumber,
				Organization: asn.AutonomousSystemOrganization,
			},
			ISP:         enterpriseISP,
			AnonymousIP: enterpriseAnonymousIP,
		},
		Travel:    travel,
		UserAgent: indata.Data.UserAgent,
	}

	if err := loginEventCurrent.AddHash(); err != nil {
		return nil, err
	}

	if _, err := c.store.AddLoginEvent(ctx, loginEventCurrent); err != nil {
		return nil, err
	}

	// decide if the login is reliable
	if loginEventLatest.Hash == loginEventCurrent.Hash{
		
	}



	return &ReplyLoginEvent{Status: true}, nil
}
