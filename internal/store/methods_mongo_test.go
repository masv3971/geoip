package store

import (
	"context"
	"fmt"
	"geoip/pkg/model"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	mockEppn     = "testUser"
)

func TestAddLoginEvent(t *testing.T) {
	if os.Getenv("backOffVscode") == "" {
		t.SkipNow()
	}

	s := mockNew(t)
	ctx := context.TODO()

	loginEvent := &model.LoginEvent{
		Timestamp: time.Now(),
		EppnHashed:      mockEppn,
		IP:        &model.IP{},
		KnownDevice:  true,
	}

	loginEventResp, err := s.Doc.AddLoginEvent(ctx, loginEvent)
	assert.NoError(t, err)

	fmt.Println("loginEventResp", loginEventResp)
}

func TestGetLoginEvents(t *testing.T) {
	if os.Getenv("backOffVscode") == "" {
		t.SkipNow()
	}

	s := mockNew(t)
	ctx := context.TODO()

	loginEvent, err := s.Doc.GetLoginEvents(ctx, mockEppn)
	assert.NoError(t, err)
	fmt.Println("loginEvent", loginEvent)
}

func TestPing(t *testing.T) {
	if os.Getenv("backOffVscode") == "" {
		t.SkipNow()
	}

	s := mockNew(t)
	ctx := context.TODO()

	err := s.Doc.ping(ctx)
	assert.NoError(t, err)
}

func TestCreateIndexes(t *testing.T) {
	if os.Getenv("backOffVscode") == "" {
		t.SkipNow()
	}

	s := mockNew(t)
	ctx := context.TODO()

	err := s.Doc.createLoginEventsIndexes(ctx)
	assert.NoError(t, err)
}

func TestGetLatestLoginEvent(t *testing.T) {
	if os.Getenv("backOffVscode") == "" {
		t.SkipNow()
	}

	s := mockNew(t)
	ctx := context.TODO()

	loginEvent1 := &model.LoginEvent{
		EppnHashed: mockEppn,
	}
	loginEvent1.UserAgent.Browser.Family = "BlackBerry9700"

	loginEvent2 := &model.LoginEvent{
		EppnHashed: mockEppn,
	}
	loginEvent2.UserAgent.Browser.Family = "Mozilla"

	_, err := s.Doc.AddLoginEvent(ctx, loginEvent1)
	assert.NoError(t, err)

	_, err = s.Doc.AddLoginEvent(ctx, loginEvent2)
	assert.NoError(t, err)

	loginEvent, err := s.Doc.GetLatestLoginEvent(ctx, mockEppn)
	assert.NoError(t, err)

	assert.Equal(t, loginEvent2.UserAgent.Browser.Family, loginEvent.UserAgent.Browser.Family)

	err = s.Doc.RemoveLoginEventForUser(ctx, mockEppn)
	assert.NoError(t, err)
}

func TestLoginEvents(t *testing.T) {
	if os.Getenv("backOffVscode") == "" {
		t.SkipNow()
	}

	s := mockNew(t)
	ctx := context.TODO()

	tts := []struct {
		name       string
		loginEvent *model.LoginEvent
	}{
		{
			name:       "OK",
			loginEvent: model.MockLoginEvent(model.MockConfig{Suffix: "a"}),
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.Doc.AddLoginEvent(ctx, tt.loginEvent)
			assert.NoError(t, err)

		})
	}
}
