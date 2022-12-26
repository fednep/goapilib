package common

import (
	"errors"
	"fmt"
	"net/http"
	"time"

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

	Timeout int
}

func (cfg *ServerConfig) LoadFromEnv() {

	// Server Configuration overrides
	config.Env("HTTP_ADDRESS", &cfg.Address)
	config.EnvInt("HTTP_PORT", &cfg.Port)

	config.EnvInt("HTTP_TIMEOUT", &cfg.Timeout)

	config.EnvBool("HTTP_USE_TLS", &cfg.UseTLS)
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

func (cfg *ServerConfig) Server(handler http.Handler) *http.Server {

	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.Address, cfg.Port),

		Handler: handler,

		ReadTimeout:  time.Duration(cfg.Timeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Timeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Timeout) * time.Second,
	}

	return srv
}
