package server

import (
	"00-newapp-template/internal/pkg"
	"00-newapp-template/internal/pkg/server/db"
	"00-newapp-template/pkg/acme/cache"
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
	Context   context.Context
	Router    chi.Router
	HTTP      *http.Server
	Finished  context.CancelFunc
	DB        db.SimpleDB
	DiskCache *cache.Disk

	Log         *log.Logger
	CacheFolder string
	ListenPort  string
}

// NewServer configs the HTTP, router, context, log and a DB to mock the ACME HTTP API
func NewServer(config *pkg.Config) (server Server) {
	server.Log = config.Log
	server.ListenPort = config.Server.ListenPort
	server.CacheFolder = config.Server.CacheFolder

	server.Context = config.Context
	server.Router = chi.NewRouter()
	server.HTTP = &http.Server{Addr: ":" + server.ListenPort, Handler: server.Router}
	server.DB = db.NewSimpleDB()
	return
}

// EnableCache will create a new Disk Cache for all request.
func (s *Server) EnableCache(cacheFolder string, cryptoKey string) {
	var useCrypto = false
	if cryptoKey != "" {
		useCrypto = true
	}
	s.DiskCache = cache.NewDisk(cacheFolder, cryptoKey, useCrypto)
	return
}

// ListenAndServe will attempt to bind and provide HTTP service. It's hooked for signals and smooth shutdown.
func (s *Server) ListenAndServe() (err error) {
	s.hookShutdownSignal()

	go func() {
		s.Log.Infof("s started")
		err = s.HTTP.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.Log.Errorf("error serving: %+v", err)
		}
		s.Finished()
	}()

	select {
	case <-s.Context.Done():
		s.Log.Infof("s stopped")
	}

	return
}
func (s *Server) hookShutdownSignal() {
	stop := make(chan os.Signal)

	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)

	s.Context, s.Finished = context.WithCancel(s.Context)
	go func() {
		sig := <-stop
		s.Log.Infof("s termination signal '%s' received", sig)
		s.Finished()
	}()

	return
}
