package config

import (
	"github.com/BurntSushi/toml"
)

func LoadToml(cfg any, fn string) error {
	_, err := toml.DecodeFile(fn, cfg)
	if err != nil {
		return err
	}

	return nil
}
