package common

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

// ServerConfig contains common config parameters
// for configuring HTTP(S) server
type ServerConfig struct {
	Address  string `env:"HTTP_ADDRESS"`
	Port     int    `env:"HTTP_PORT"`
	UseTLS   bool   `env:"HTTP_USE_TLS"`
	CertFile string `env:"HTTP_CERT_FILE"`
	KeyFile  string `env:"HTTP_KEY_FILE"`

	Timeout int `env:"HTTP_TIMEOUT"`
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
