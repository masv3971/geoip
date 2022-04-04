package rules

import (
	"geoip/pkg/logger"
	"geoip/pkg/model"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockContries = []byte(`
---

countries:
  - Afghanistan: 10.0
  - Sweden: 0.01
`)

func mockClient(t *testing.T) *Client {
	tempDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tempDir, "countries.yaml"), mockContries, fs.FileMode(os.O_RDWR))
	assert.NoError(t, err)

	cfg := &model.Cfg{
		Rules: struct {
			Folder string "yaml:\"folder\""
		}{
			Folder: tempDir,
		},
	}
	c, err := New(cfg, logger.New("test", false).New("rule"))
	assert.NoError(t, err)

	return c
}

func TestLoadRules(t *testing.T) {
	c := mockClient(t)

	c.loadRules()
}
