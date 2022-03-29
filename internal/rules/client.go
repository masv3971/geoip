package rules

import (
	"geoip/pkg/logger"
	"geoip/pkg/model"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Client holds rule object
type Client struct {
	log   *logger.Logger
	cfg   *model.Cfg
	rules *model.Rules
}

// New creates a new instance of rule
func New(cfg *model.Cfg, log *logger.Logger) (*Client, error) {
	c := &Client{
		log: log,
		cfg: cfg,
	}

	if err := c.loadRules(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) loadRules() error {
	err := filepath.Walk("rules",
		func(path string, info fs.FileInfo, err error) error {
			if filepath.Ext(info.Name()) != "yaml" {
				return nil
			}
			if err != nil {
				return err
			}
			fileData, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			c.rules = &model.Rules{}
			yaml.Unmarshal(fileData, c.rules)
			return nil
		})

	if err != nil {
		return err
	}
	return nil
}
