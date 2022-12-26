package endpoint

import (
	"log"
	"net/http"

	"github.com/fednep/goapilib/config/common"
)

// ServeForever starts HTTP server
// returns on error or after server is terminated
func ServeForever(cfg common.ServerConfig, handler http.Handler) error {

	server := cfg.Server(handler)

	if cfg.UseTLS {

		log.Printf("Starting server with TLS on: %s", server.Addr)

		return http.ListenAndServeTLS(
			server.Addr,
			cfg.CertFile,
			cfg.KeyFile,
			handler)
	}

	log.Printf("Starting server on: %s", server.Addr)
	return server.ListenAndServe()
}
