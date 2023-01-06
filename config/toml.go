package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// LoadConfig fill cfg struct from multiple places:
// 1. config file in TOML format if -config args is specified
// 3. environment variables
//
// Parameters:
// cfg - config struct
// loadOverrides - load variables from environment vars
func LoadToml(cfg any, fn string) error {
	_, err := toml.DecodeFile(fn, &cfg)
	if err != nil {
		return fmt.Errorf("config load error: %w", err)
	}

	return nil
}
