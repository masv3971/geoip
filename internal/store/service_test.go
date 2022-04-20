package store

import (
	"context"
	"geoip/pkg/logger"
	"geoip/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockNew(t *testing.T) *Service {
	cfg := &model.Cfg{}
	cfg.Mongodb.Addr = "mongodb://127.0.0.1:27017"

	s, err := New(context.TODO(), cfg, logger.New("test", false).New("test"))
	assert.NoError(t, err)

	return s
}
