package traveler

import (
	"context"
	"geoip/pkg/model"
)

type Client struct {
	cfg   *model.Cfg
	store storage
}

type storage interface {
	GetLatestLoginEvent(ctx context.Context, eppn string) (*model.LoginEvent, error)
}

// New create a new instance of traveler client
func New(ctx context.Context) (*Client, error) {
	c := &Client{}

	return c, nil
}
