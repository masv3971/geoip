package httpserver

import (
	"context"
	"geoip/pkg/helpers"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Service) middlewareDuration(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		duration := time.Since(t)
		c.Set("duration", duration)
	}
}

func (s *Service) middlewareTraceID(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("sunet-request-id", uuid.NewString())
		c.Header("sunet-request-id", c.GetString("sunet-request-id"))
		c.Next()
	}
}

func (s *Service) middlewareSolidContentType(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodHead:
			c.Next()
		case http.MethodGet:
			c.Next()
		default:
			if c.Request.Method != http.MethodGet {
				if c.ContentType() == "" {
					renderContent(c, 400, gin.H{"message": "The given input is not supported by the server configuration."})
				}
			}
			c.Next()
		}
	}
}

func (s *Service) middlewareLogger(ctx context.Context) gin.HandlerFunc {
	log := s.logger.New("http")
	return func(c *gin.Context) {
		c.Next()
		log.Info("request", "status", c.Writer.Status(), "url", c.Request.URL.String(), "method", c.Request.Method, "sunet-request-id", c.GetString("sunet-request-id"))
	}
}

func (s *Service) middlewareCrash(ctx context.Context) gin.HandlerFunc {
	log := s.logger.New("http")
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				status := c.Writer.Status()
				log.Error("crash", "error", r, "status", status, "url", c.Request.URL.Path, "method", c.Request.Method)
				renderContent(c, 500, gin.H{"data": nil, "error": helpers.NewError("internal_server_error")})
			}
		}()
		c.Next()
	}
}
