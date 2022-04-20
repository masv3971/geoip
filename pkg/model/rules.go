package model

// Rules holds all the rule definitions
type Rules struct {
	Countries RulesCountries `yaml:"countries"`
}

// RulesCountries countries rule definition
type RulesCountries map[string]float64
