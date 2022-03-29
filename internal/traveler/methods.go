package traveler

import (
	"context"
	"geoip/pkg/model"
	"math"
	"time"
)

// Travel tries to travel really fast, if the duration is greater than 930km/h it's too fast to be possible
func (c *Client) Travel(ctx context.Context, previous, current *model.LoginEvent) (*model.Travel, error) {
	if previous == nil || current == nil {
		return nil, nil
	}
	hsin := func(theta float64) float64 {
		return math.Pow(math.Sin(theta/2), 2)
	}

	// Convert to from degree to rand
	latCurrent := current.Location.Coordinates.Latitude * math.Pi / 180
	longCurrent := current.Location.Coordinates.Longitude * math.Pi / 180
	latLatest := previous.Location.Coordinates.Latitude * math.Pi / 180
	longLatest := previous.Location.Coordinates.Longitude * math.Pi / 180

	earthRadius := float64(6378100) //meters

	h := hsin(latLatest-latCurrent) + math.Cos(latCurrent)*math.Cos(latLatest)*hsin(longLatest-longCurrent)

	distance := 2 * earthRadius * math.Asin(math.Sqrt(h))

	travelDuration := time.Now().Sub(previous.Timestamp)

	fastAirplane := 258.333333 // meter/second

	travelSpeed := distance / travelDuration.Seconds()

	impssibleTravel := false
	if travelSpeed > fastAirplane {
		impssibleTravel = true
	}

	travel := &model.Travel{
		Distance:           distance,
		DistanceUnit:       "meter",
		TravelDuration:     time.Duration(travelDuration.Seconds()),
		TravelDurationUnit: "second",
		TravelSpeed:        travelSpeed,
		TravelSpeedUnit:    "meter/second",
		IsTravelImpossible: impssibleTravel,
	}

	return travel, nil
}
