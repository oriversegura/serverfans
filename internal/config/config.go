package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	IP       string `json:"ip"`
	User     string `json:"user"`
	Password string `json:"password"`
	Speed    string `json:"speed"`
}

// Load reads a JSON config file if one is passed as a CLI argument.
// Missing or unreadable files are silently ignored — the UI will prompt for missing values.
func Load() *Config {
	cfg := &Config{}
	if len(os.Args) < 2 {
		return cfg
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		return cfg
	}
	json.Unmarshal(data, cfg) //nolint:errcheck — partial config is fine
	return cfg
}
