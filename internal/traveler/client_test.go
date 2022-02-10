package traveler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockNew(t *testing.T) *Client {
	c, err := New(context.TODO())
	assert.NoError(t, err)

	return c
}
