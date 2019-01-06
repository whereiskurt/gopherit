package server

import (
	"00-newapp-template/pkg/config"
	"00-newapp-template/pkg/metrics"
	"00-newapp-template/pkg/server"
	log "github.com/sirupsen/logrus"
	"time"
)

// Start will create a new Server, attach a Handler, and start listening on the port logging to the log.
func Start(config *config.Config, metrics *metrics.Metrics) {
	config.Server.EnableLogging()

	config.Log.Debugf("internal.app.cmd.server.Start called with -> config(%+v) and metrics->(%+v)", config, metrics)
	l := config.Log.WithFields(log.Fields{
		"docroot": config.Server.RootFolder,
		"cache":   config.Server.CacheFolder,
		"port":    config.Server.ListenPort,
	})

	s := server.NewServer(config, metrics)
	s.EnableDefaultRouter()

	go func() {
		for {
			time.Sleep(15 * time.Minute)
			config.Server.DumpMetrics()
		}
	}()

	l.Info("starting server")
	_ = s.ListenAndServe()
	l.Info("server stopped.")

	l.Info("dumping metrics for server")
	config.Server.DumpMetrics()

	return
}
