package traveler

import (
	"context"
	"geoip/pkg/model"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTravel(t *testing.T) {
	type args struct {
		lat, long  float64
		loginEvent *model.LoginEvent
	}
	tts := []struct {
		name string
		have args
		want *model.Travel
	}{
		{
			name: "zero distance",
			have: args{
				lat:  59.3274,
				long: 18.0653,
				loginEvent: &model.LoginEvent{
					IP: &model.IP{
						Coordinates: &model.Coordinates{
							Latitude:  59.3274,
							Longitude: 18.0653,
						},
					},
					TimeStamp: time.Now().Add(-24 * time.Hour),
				},
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
			have: args{
				lat:  37.751,
				long: -97.822,
				loginEvent: &model.LoginEvent{
					IP: &model.IP{
						Coordinates: &model.Coordinates{
							Latitude:  59.3274,
							Longitude: 18.0653,
						},
					},
					TimeStamp: time.Now().Add(-24 * time.Hour),
				},
			},
			want: &model.Travel{
				Distance:           7.734844247101788e+06,
				DistanceUnit:       "meter",
				TravelDuration:     86400,
				TravelDurationUnit: "second",
				TravelSpeed:        89.52366009179627,
				TravelSpeedUnit:    "meter/second",
				IsTravelImpossible: true,
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			s := mockNew(t)

			got, err := s.Travel(context.TODO(), tt.have.lat, tt.have.long, tt.have.loginEvent)
			assert.NoError(t, err)

			assert.Equal(t, tt.want.Distance, got.Distance)
			assert.Equal(t, tt.want.IsTravelImpossible, got.IsTravelImpossible)
			assert.Equal(t, math.Round(tt.want.TravelSpeed*100)/100, math.Round(got.TravelSpeed*100)/100)
		})
	}
}
