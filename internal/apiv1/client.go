package apiv1

import (
	"context"
	"geoip/internal/maxmind"
	"geoip/internal/store"
	"geoip/internal/traveler"
	"geoip/pkg/logger"
	"geoip/pkg/model"
)

// Client holds apiv1 object
type Client struct {
	config   *model.Cfg
	logger   *logger.Logger
	maxmind  *maxmind.Service
	traveler *traveler.Client
	store    storage
}

type storage interface {
	AddLoginEvent(ctx context.Context, loginEvent *model.LoginEvent) (interface{}, error)
	GetLatestLoginEvent(ctx context.Context, eppn string) (*model.LoginEvent, error)
}

// New creates a new instance of apiv1
func New(ctx context.Context, cfg *model.Cfg, maxmind *maxmind.Service, store *store.Service, traveler *traveler.Client, log *logger.Logger) (*Client, error) {
	c := &Client{
		config:   cfg,
		logger:   log,
		maxmind:  maxmind,
		store:    store.Doc,
		traveler: traveler,
	}

	return c, nil
}
