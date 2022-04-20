package httpserver

import (
	"context"
	"geoip/internal/apiv1"
)

// Apiv1 interface
type Apiv1 interface {
	HandlerLoginEvent(ctx context.Context, indata *apiv1.RequestLoginEvent) (*apiv1.ReplyLoginEvent, error)

	HandlerStatsOverview(ctx context.Context, indata *apiv1.RequestStatsOverview) (*apiv1.ReplyStatsOverview, error)
	HandlerStatsEppnLong(ctx context.Context, indata *apiv1.RequestStatsEppnLong) (*apiv1.ReplyStatsEppnLong, error)
	HandlerStatsEppnSpecific(ctx context.Context, indata *apiv1.RequestStatsEppnSpecific) (*apiv1.ReplyStatsEppnSpecific, error)
}
