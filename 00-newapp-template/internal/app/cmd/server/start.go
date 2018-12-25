package server

import (
	"00-newapp-template/internal/pkg"
	"00-newapp-template/internal/pkg/server"
)

// Start will create a new Server, attach a Router, and start listening on the port logging to the log.
func Start(config *pkg.Config) {
	s := server.NewServer(config)
	s.NewRouter()
	_ = s.ListenAndServe()

	return
}
