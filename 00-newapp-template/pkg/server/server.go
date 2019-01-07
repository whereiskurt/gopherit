package server

import (
	"00-newapp-template/pkg/acme"
	"00-newapp-template/pkg/cache"
	"00-newapp-template/pkg/config"
	"00-newapp-template/pkg/metrics"
	"00-newapp-template/pkg/server/db"
	"00-newapp-template/pkg/server/middleware"
	"context"
	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// Server is built on go-chi
type Server struct {
	Context           context.Context
	Handler           chi.Router
	HTTP              *http.Server
	Finished          context.CancelFunc
	DB                db.SimpleDB
	DiskCache         *cache.Disk
	Log               *log.Logger
	CacheFolder       string
	ListenPort        string
	Metrics           *metrics.Metrics
	MetricsListenPort string
}

// NewServer configs the HTTP, router, context, log and a DB to mock the ACME HTTP API
func NewServer(config *config.Config, metrics *metrics.Metrics) (server Server) {
	server.Log = config.Log
	server.ListenPort = config.Server.ListenPort
	server.CacheFolder = config.Server.CacheFolder
	server.MetricsListenPort = config.Server.MetricsListenPort

	if config.Server.CacheResponse {
		server.EnableCache(config.Server.CacheFolder, config.Server.CacheKey)
	}

	server.Context = config.Context
	server.Handler = chi.NewRouter()

	server.HTTP = &http.Server{
		Addr: ":"+server.ListenPort,
		Handler: server.Handler,
		// Ian Kent recommends these timeouts be set:
		//   https://www.youtube.com/watch?v=YF1qSfkDGAQ&t=333s
		IdleTimeout:time.Duration(time.Second), // This maybe tooooo aggressive..*shrugs* :)
		ReadTimeout:time.Duration(time.Second),
		WriteTimeout:time.Duration(time.Second),
	}

	server.DB = db.NewSimpleDB()
	server.Metrics = metrics
	return
}

// ListenAndServe will attempt to bind and provide HTTP service. It's hooked for signals and smooth shutdown.
func (s *Server) ListenAndServe() (err error) {
	s.hookShutdownSignal()

	// Start the /metrics server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		_ = http.ListenAndServe(":"+s.MetricsListenPort, nil)
	}()

	// Start the HTTP server
	go func() {
		s.Log.Infof("server started")
		err = s.HTTP.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.Log.Errorf("error serving: %+v", err)
		}
		s.Finished()
	}()

	select {
	case <-s.Context.Done():
		s.Log.Infof("server stopped")
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
		s.Log.Infof("termination signal '%s' received for server", sig)
		s.Finished()
	}()

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

func (s *Server) cacheClear(r *http.Request, endPoint acme.EndPointType, service metrics.EndPointType) {
	if s.DiskCache == nil {
		return
	}
	if s.Metrics != nil {
		s.Metrics.CacheInc(service, metrics.Methods.Cache.Invalidate)
	}

	filename, _ := acme.ToCacheFilename(endPoint, middleware.ContextMap(r))
	filename = filepath.Join(".", s.DiskCache.CacheFolder, filename)

	s.DiskCache.Clear(filename)
}
func (s *Server) cacheStore(r *http.Request, w http.ResponseWriter, endPoint acme.EndPointType, service metrics.EndPointType, bb []byte) {
	if s.DiskCache == nil {
		return
	}
	// Metrics!
	if s.Metrics != nil {
		s.Metrics.CacheInc(service, metrics.Methods.Cache.Store)
	}

	filename, _ := acme.ToCacheFilename(endPoint, middleware.ContextMap(r))
	prettyCache := middleware.NewPrettyPrint(w).Prettify(bb)
	_ = s.DiskCache.Store(filename, prettyCache)
}
func (s *Server) cacheFetch(r *http.Request, endPoint acme.EndPointType, service metrics.EndPointType) (bb []byte, err error) {
	if s.DiskCache == nil {
		return
	}

	filename, _ := acme.ToCacheFilename(endPoint, middleware.ContextMap(r))
	filename = filepath.Join(".", s.DiskCache.CacheFolder, filename)

	bb, err = s.DiskCache.Fetch(filename)

	if err == nil && len(bb) > 0 && s.Metrics != nil {
		s.Metrics.CacheInc(service, metrics.Methods.Cache.Hit)
	} else {
		s.Metrics.CacheInc(service, metrics.Methods.Cache.Miss)
	}

	return
}
