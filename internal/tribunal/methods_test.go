package tribunal

import (
	"geoip/pkg/model"
	"testing"
)

const (
	THash1 = "98023412311609234234"
	THash2 = "11782566312679053365"
)

func TestEstimateFish(t *testing.T) {
	tts := []struct {
		name string
		have *Client
		want *Verdict
	}{
		{
			name: "No previous login",
			have: &Client{
				Previous: nil,
				Current: &model.LoginEvent{
					Eppn: "testUser",
					Hash: THash1,
				},
			},
			want: &Verdict{
				Reason:    ReasonNoPreviousLogin,
				Sentence:  SentenceNone,
				FishScore: 0,
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
