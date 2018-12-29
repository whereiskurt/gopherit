package server

import (
	"00-newapp-template/internal/pkg"
	"00-newapp-template/internal/pkg/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// Start will create a new Server, attach a Router, and start listening on the port logging to the log.
func Start(config *pkg.Config, metrics *pkg.Metrics) {
	s := server.NewServer(config, metrics)
	s.NewRouter()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":"+config.Metrics.ListenPort, nil)
	}()

	_ = s.ListenAndServe()

	return
}
