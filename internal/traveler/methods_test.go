package traveler

import (
	"context"
	"geoip/pkg/model"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTravel(t *testing.T) {
	type args struct {
		lat, long  float64
		loginEvent *model.LoginEvent
	}
	type have struct {
		current  *model.LoginEvent
		previous *model.LoginEvent
	}
	tts := []struct {
		name string
		have have
		want *model.Travel
	}{
		{
			name: "zero distance",
			have: have{
				current:  model.MockLoginEvent(model.MockConfig{Country: "sweden", H: 14}),
				previous: model.MockLoginEvent(model.MockConfig{Country: "sweden"}),
			},
			want: &model.Travel{
				Distance:           0,
				DistanceUnit:       "meter",
				TravelDuration:     86400,
				TravelDurationUnit: "second",
				TravelSpeed:        0,
				TravelSpeedUnit:    "meter/second",
				IsTravelImpossible: false,
			},
		},
		{
			name: "stockholm->USA",
			have: have{
				current:  model.MockLoginEvent(model.MockConfig{Country: "sweden", H: 17}),
				previous: model.MockLoginEvent(model.MockConfig{Country: "usa", H: 12}),
			},
			want: &model.Travel{
				Distance:           6.584876652588911e+06,
				DistanceUnit:       "meter",
				TravelDuration:     86400,
				TravelDurationUnit: "second",
				TravelSpeed:        365.83,
				TravelSpeedUnit:    "meter/second",
				IsTravelImpossible: true,
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			s := mockNew(t)

			got, err := s.Travel(context.TODO(), tt.have.previous, tt.have.current)
			assert.NoError(t, err)

			assert.Equal(t, tt.want.Distance, got.Distance, "distance")
			assert.Equal(t, tt.want.IsTravelImpossible, got.IsTravelImpossible, "isTravelImpossible")
			assert.Equal(t, math.Round(tt.want.TravelSpeed*100)/100, math.Round(got.TravelSpeed*100)/100, "travelSpeed")
		})
	}
}
