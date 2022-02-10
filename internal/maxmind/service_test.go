package maxmind

import (
	"context"
	"geoip/pkg/logger"
	"geoip/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockNew(t *testing.T) *Service {
	cfg := &model.Cfg{}

	service, err := New(context.TODO(), cfg, nil, logger.New("test", false).New("test"))
	assert.NoError(t, err)

	return service
}
