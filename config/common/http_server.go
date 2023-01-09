package common

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fednep/goapilib/config"
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

func (cfg ServerConfig) IsValid() error {

	if cfg.Port == 0 {
		return config.FieldError{FieldName: "Port", Message: "not configured"}
	}

	if cfg.UseTLS {
		if cfg.CertFile == "" {
			return config.FieldError{
				FieldName: "CertFile",
				Message:   "cannot be empty when TLS is enabled"}
		}

		if cfg.KeyFile == "" {
			return config.FieldError{
				FieldName: "KeyFile",
				Message:   "cannot be empty when TLS is enabled"}
		}
	}

	return nil
}

func (cfg ServerConfig) Server(handler http.Handler) *http.Server {

	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.Address, cfg.Port),

		Handler: handler,

		ReadTimeout:  time.Duration(cfg.Timeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Timeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Timeout) * time.Second,
	}

	return srv
}
