package statistics

import (
	"geoip/pkg/model"

	"github.com/montanaflynn/stats"
)

// StatsData holds data
type StatsData map[string]*model.StatsData

// NewSpecific return
func NewSpecific(loginEvents model.LoginEvents) (StatsData, error) {

	ips := []string{}
	uaOS := []int{}
	uaDevice := []int{}
	uaBrowser := []int{}

	for _, loginEvent := range loginEvents {
		ips = append(ips, loginEvent.IP.IPAddrML)
		uaOS = append(uaOS, loginEvent.UserAgent.OS.FamilyML)
		uaDevice = append(uaDevice, loginEvent.UserAgent.Device.FamilyML)
		uaBrowser = append(uaBrowser, loginEvent.UserAgent.Browser.FamilyML)
	}
	mura := map[string]interface{}{
		"ip":                ips,
		"user_agent_device": uaDevice,
		"user_agent_os":     uaOS,
	}

	sData := StatsData{}

	for name, value := range mura {
		s := &model.StatsData{}
		var err error
		stdFormat := stats.LoadRawData(value)

		s.Len = stdFormat.Len()

		s.StandardDeviation, err = stdFormat.StandardDeviation()
		if err != nil {
			return nil, err
		}

		sData[name] = s

	}

	return sData, nil

}
