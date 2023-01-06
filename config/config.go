package config

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// Section defines interface used by configuration sections
type Section interface {
	IsValid() error
}

func IsValid(config any) error {
	return nil
}

// Load fills config structure based on the default rules
//
// These default rules are:
//  1. look for config.toml file, decodes to a config if exists
//  2. look for .env file, load it as environment vars
//  3. look up for any env variables correspondong to
//     config fields and load them
//  4. validate config
func Load(config any) error {
	return LoadConfig(config, DefaultRules)
}

func LoadConfig(config any, rules Rules) error {

	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	tomlFileName := &rules.TomlFile
	envFileName := &rules.DotEnvFile

	if rules.TomlArgsOption != "" {
		tomlFileName = flagSet.String(
			rules.TomlArgsOption,
			rules.TomlFile,
			"Specify config file in TOML format")
	}

	if rules.DotEnvArgsOption != "" {
		envFileName = flagSet.String(
			rules.DotEnvArgsOption,
			rules.DotEnvFile,
			"Specify config file in .env (key=value) format")
	}

	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		return fmt.Errorf("cannot parse flags: %w", err)
	}

	if rules.Toml {
		if err := loadToml(config, *tomlFileName, rules.TomlRequired); err != nil {
			return err
		}
	}

	if rules.DotEnv {
		if err := loadDotEnv(*envFileName, rules.DotEnvRequired); err != nil {
			return err
		}
	}

	return nil
}

func loadToml(cfg any, fn string, isRequired bool) error {

	fnExists, err := checkRequired(fn, isRequired)
	if err != nil {
		return err
	}

	if fnExists {
		log.Printf("Loading config from %q", fn)
		err := LoadToml(cfg, fn)
		if err != nil {
			return err
		}
	}

	return nil
}

func loadDotEnv(fn string, isRequired bool) error {

	fnExists, err := checkRequired(fn, isRequired)
	if err != nil {
		return err
	}

	if fnExists {
		log.Printf("Loading env file %q", fn)
		err := LoadEnvFile(fn)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkRequired(fn string, isRequired bool) (bool, error) {
	if fn == "" && isRequired {
		return false, fmt.Errorf("config file not specified")
	}

	fnExists, err := fileExists(fn)
	if err != nil {
		return false, fmt.Errorf("error accessing %s: %w", fn, err)
	}

	if !fnExists && isRequired {
		return false, fmt.Errorf("config file %q not found", fn)
	}

	return fnExists, nil
}

func fileExists(fn string) (bool, error) {
	_, err := os.Stat(fn)
	if err == nil {
		return true, nil
	}

	if err == os.ErrNotExist {
		return false, nil
	}

	return false, err
}
