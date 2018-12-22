package server

import (
	"00-newapp-template/internal/pkg/server/db"
	"context"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Server is built on go-chi
type Server struct {
	Context  context.Context
	Router   chi.Router
	HTTP     *http.Server
	Log      *log.Logger
	Finished context.CancelFunc
	DB       db.SimpleDB
}

// NewServer configs the HTTP, router, context, log and a DB to mock the ACME HTTP API
func NewServer(context context.Context, listenPort string, log *log.Logger) (server Server) {
	server.Context = context
	server.Router = chi.NewRouter()
	server.HTTP = &http.Server{Addr: ":" + listenPort, Handler: server.Router}
	server.Log = log
	server.DB = db.NewSimpleDB()
	return
}

// ListenAndServe will attempt to bind and provide HTTP service. It's hooked for signals and smooth shutdown.
func (server *Server) ListenAndServe() (err error) {
	server.hookShutdownSignal()

	go func() {
		server.Log.Infof("server starting")
		err = server.HTTP.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			server.Log.Errorf("error serving: %+v", err)
		}
		server.Finished()
	}()

	select {
	case <-server.Context.Done():
		server.Log.Infof("server stopped")
	}

	return
}
func (server *Server) hookShutdownSignal() {
	stop := make(chan os.Signal)

	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)

	server.Context, server.Finished = context.WithCancel(server.Context)
	go func() {
		sig := <-stop
		server.Log.Infof("server termination signal '%s' received", sig)
		server.Finished()
	}()

	return
}
