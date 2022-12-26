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
		log.Printf("Starting HTTPS server on: %s", server.Addr)
		err := http.ListenAndServeTLS(
			server.Addr,
			cfg.CertFile,
			cfg.KeyFile,
			handler)
		log.Printf("HTTPS server (%s) stopped: %s", server.Addr, err)
		return err
	}

	log.Printf("Starting HTTP server on: %s", server.Addr)
	err := server.ListenAndServe()
	log.Printf("HTTP Server (%s) stopped: %s", server.Addr, err)
	return err
}
