package common

import (
	"errors"

	"github.com/fednep/goapilib/config"
)

// ServerConfig contains common config parameters
// for configuring HTTP(S) server
type ServerConfig struct {
	Address  string
	Port     int
	UseTLS   bool
	CertFile string
	KeyFile  string

	ReadTimeoutSec  int
	WriteTimeoutSec int
}

func (cfg *ServerConfig) LoadEnvOverrides() {

	// Server Configuration overrides
	config.Env("HTTP_ADDRESS", &cfg.Address)

	// If specified any of this, we'll use it (for backward compatibility)
	config.EnvInt("HTTP_PORT", &cfg.Port)

	config.EnvBool("HTTP_USE_TLS", &cfg.UseTLS)

	config.EnvInt("HTTP_READ_TIMEOUT", &cfg.ReadTimeoutSec)
	config.EnvInt("HTTP_WRITE_TIMEOUT", &cfg.WriteTimeoutSec)

	if cfg.UseTLS {
		config.Env("HTTP_CERT_FILE", &cfg.CertFile)
		config.Env("HTTP_KEY_FILE", &cfg.KeyFile)
	}
}

func (cfg *ServerConfig) IsValid() error {

	if cfg.Port == 0 {
		return errors.New("[server] need to have port configured")
	}

	if cfg.UseTLS {
		if cfg.CertFile == "" || cfg.KeyFile == "" {
			return errors.New("[server] set to use TLS but does not have 'certfile' or 'keyfile' configured")
		}
	}

	return nil
}
