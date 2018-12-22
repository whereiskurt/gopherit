package server

import (
	"00-newapp-template/internal/pkg/server"
	"context"
	log "github.com/sirupsen/logrus"
)

// Start will create a new Server, attach a Router, and start listening on the port logging to the log.
func Start(context context.Context, port string, log *log.Logger) {
	s := server.NewServer(context, port, log)
	s.NewRouter()
	_ = s.ListenAndServe()

	return
}
