package httpserver

import (
	"context"
	"geoip/internal/apiv1"

	"github.com/gin-gonic/gin"
)

func (s *Service) endpointLoginEvent(ctx context.Context, c *gin.Context) (interface{}, error) {
	request := &apiv1.RequestLoginEvent{}
	if err := s.bindRequest(c, request); err != nil {
		return nil, err
	}
	reply, err := s.apiv1.HandlerLoginEvent(ctx, request)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (s *Service) endpointStatsOverview(ctx context.Context, c *gin.Context) (interface{}, error) {
	request := &apiv1.RequestStatsOverview{}
	if err := s.bindRequest(c, request); err != nil {
		return nil, err
	}
	reply, err := s.apiv1.HandlerStatsOverview(ctx, request)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
