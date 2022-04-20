package model

import "time"

type Cfg struct {
	Production bool `yaml:"production"`
	APIServer  struct {
		Addr string `yaml:"addr"`
	} `yaml:"api_server"`
	Rules struct {
		Folder string `yaml:"folder"`
	} `yaml:"rules"`
	MaxMind struct {
		UpdatePeriodicity time.Duration `yaml:"update_periodicity"`
		LicenseKey        string        `yaml:"license_key"`
		Enterprise        bool          `yaml:"enterprise"`
	} `yaml:"maxmind"`
	Mongodb struct {
		Addr string `yaml:"addr" validate:"required"`
	} `yaml:"mongodb"`
	KVStorage struct {
		Redis struct {
			DB                  int      `yaml:"db" validate:"required"`
			Addr                string   `yaml:"addr" validate:"required_without_all=SentinelHosts SentinelServiceName"`
			SentinelHosts       []string `yaml:"sentinel_hosts" validate:"required_without=Addr,omitempty,min=2,max=4"`
			SentinelServiceName string   `yaml:"sentinel_service_name" validate:"required_with=SentinelHosts"`
		} `yaml:"redis"`
	} `yaml:"kv_storage"`
}

// Config represent the complete config file structure
type Config struct {
	GeoIP Cfg `yaml:"geoip"`
}
