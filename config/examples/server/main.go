package main

import (
	"flag"
	"log"
	"os"

	"github.com/fednep/goapilib/config"
	"github.com/fednep/goapilib/config/common"
)

type Config struct {
	Server common.ServerConfig `env:"SERVER"`
}

func loadConfig() (*Config, error) {

	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	tomlFile := flagSet.String("config", "config1.toml", "Specify config file in TOML format")
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	cfg := Config{}

	if err := config.LoadToml(&cfg, *tomlFile); err != nil {
		return nil, err
	}

	if err := config.LoadOverrides(&cfg); err != nil {
		return nil, err
	}

	if err := config.IsValid(cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func main() {

	cfg, err := loadConfig()
	if err != nil {
		log.Printf("Error loading config: %s", err)
		os.Exit(1)
	}

	log.Printf("Config: %#v", cfg)
}
