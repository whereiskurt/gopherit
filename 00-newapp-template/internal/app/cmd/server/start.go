package server

import (
	"00-newapp-template/internal/pkg"
	"00-newapp-template/internal/pkg/metrics"
	"00-newapp-template/internal/pkg/server"
)

// Start will create a new Server, attach a Router, and start listening on the port logging to the log.
func Start(config *pkg.Config, metrics *metrics.Metrics) {
	s := server.NewServer(config, metrics)
	s.NewRouter()
	_ = s.ListenAndServe()

	return
}
