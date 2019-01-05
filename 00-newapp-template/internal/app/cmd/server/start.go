package server

import (
	"00-newapp-template/pkg"
	"00-newapp-template/pkg/metrics"
	"00-newapp-template/pkg/server"
	log "github.com/sirupsen/logrus"
	"time"
)

// Start will create a new Server, attach a Router, and start listening on the port logging to the log.
func Start(config *pkg.Config, metrics *metrics.Metrics) {
	config.Server.EnableLogging()

	l := config.Log.WithFields(log.Fields{
		"docroot": config.Server.RootFolder,
		"cache":   config.Server.CacheFolder,
		"port":    config.Server.ListenPort,
	})

	l.Info("starting server")

	s := server.NewServer(config, metrics)
	s.NewRouter()

	go func() {
		for {
			time.Sleep(15 * time.Minute)
			config.Server.DumpMetrics()
		}
	}()

	_ = s.ListenAndServe()

	l.Info("server stopped.")

	l.Info("dumping metrics for server")
	config.Server.DumpMetrics()

	return
}
