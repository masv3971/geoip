package statistics

import "geoip/pkg/logger"

// Client holds statistics object
type Client struct {
	log *logger.Logger
}

// New creates a new instance of statistics
func New(log *logger.Logger) (*Client, error) {
	c := &Client{
		log: log,
	}

	return c, nil
}
