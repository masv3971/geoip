package httpserver

import (
	"context"
	"geoip/internal/apiv1"
	"geoip/pkg/helpers"
	"geoip/pkg/logger"
	"geoip/pkg/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
	"go.opentelemetry.io/otel/trace"
)

// Service is the service object for httpserver
type Service struct {
	config *model.Cfg
	logger *logger.Logger
	server *http.Server
	apiv1  Apiv1
	gin    *gin.Engine
	tp     trace.Tracer
}

// New creates a new httpserver
func New(ctx context.Context, cfg *model.Cfg, apiv1 *apiv1.Client, log *logger.Logger) (*Service, error) {
	s := &Service{
		config: cfg,
		logger: log,
		apiv1:  apiv1,
		server: &http.Server{Addr: cfg.APIServer.Addr},
	}

	switch s.config.Production {
	case true:
		gin.SetMode(gin.ReleaseMode)
	case false:
		gin.SetMode(gin.DebugMode)
	}

	apiValidator := validator.New()
	binding.Validator = &defaultValidator{
		Validate: apiValidator,
	}

	s.gin = gin.New()
	s.server.Handler = s.gin
	s.server.ReadTimeout = time.Second * 5
	s.server.WriteTimeout = time.Second * 30
	s.server.IdleTimeout = time.Second * 90

	// Middlewares
	s.gin.Use(s.middlewareTraceID(ctx))
	s.gin.Use(s.middlewareDuration(ctx))
	s.gin.Use(s.middlewareLogger(ctx))
	s.gin.Use(s.middlewareCrash(ctx))
	s.gin.Use(s.middlewareSolidContentType(ctx))
	s.gin.NoRoute(func(c *gin.Context) {
		status := http.StatusNotFound
		p := helpers.Problem404()
		c.JSON(status, gin.H{"error": p, "data": nil})
	})

	s.regEndpoint(ctx, "POST", "/api/v1/login_event", s.endpointLoginEvent)

	// Run http server
	go func() {
		err := s.server.ListenAndServe()
		if err != nil {
			s.logger.New("http").Fatal("listen_error", "error", err)
		}
	}()

	s.logger.Info("started")

	return s, nil
}

func (s *Service) regEndpoint(ctx context.Context, method, path string, handler func(context.Context, *gin.Context) (interface{}, error)) {
	s.gin.Handle(method, path, func(c *gin.Context) {
		res, err := handler(ctx, c)
		status := 200

		if err != nil {
			status = 400
		}
		renderContent(c, status, gin.H{"data": res, "error": helpers.NewErrorFromError(err)})
	})
}

func renderContent(c *gin.Context, code int, data interface{}) {
	switch c.NegotiateFormat(gin.MIMEJSON, "*/*") {
	case gin.MIMEJSON:
		c.JSON(code, data)
	case "*/*": // curl
		c.JSON(code, data)
	default:
		c.JSON(406, gin.H{"data": nil, "error": helpers.NewErrorDetails("not_acceptable", "Accept header is invalid. It should be \"application/json\".")})
	}
}

// Close closes httpserver
func (s *Service) Close(ctx context.Context) error {
	ctx.Done()
	s.logger.Info("Quit")
	return nil
}
