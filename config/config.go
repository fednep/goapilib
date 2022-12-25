package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// ConfigSection defines interface used by configuration sections
type ConfigSection interface {
	IsValid() error
	LoadEnvOverrides()
}

// LoadOverride load config parameters override from .env file
// and validates the configuration options
func LoadOverride(configs []ConfigSection) error {

	for _, section := range configs {
		section.LoadEnvOverrides()
		if err := section.IsValid(); err != nil {
			return err
		}
	}
	return nil
}

func FromEnv(key string, val *string) {
	if value, ok := os.LookupEnv(key); ok {
		*val = value
	}
}

func FromEnvInt(key string, val *int) {
	if value, ok := os.LookupEnv(key); ok {
		if v, err := strconv.Atoi(value); err == nil {
			*val = v
		}
	}
}

func FromEnvInt64(key string, val *int64) {
	if value, ok := os.LookupEnv(key); ok {
		if v, err := strconv.ParseInt(value, 10, 64); err == nil {
			*val = v
		}
	}
}

func FromEnvBool(key string, val *bool) {
	if value, ok := os.LookupEnv(key); ok {
		if v, err := strconv.ParseBool(value); err == nil {
			*val = v
		}
	}
}

func FromEnvList(key string, val *[]string) {
	if value, ok := os.LookupEnv(key); ok {
		*val = strings.Split(value, ",")
	}
}

func LoadDotEnv() error {
	// Retrieve current directory
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return fmt.Errorf("Cannot get absolute path of executable: %s", err)
	}

	LoadEnvFile(path.Join(dir, ".env"))

	return nil
}
