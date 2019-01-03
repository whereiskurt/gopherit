package server

import (
	"00-newapp-template/internal/pkg/metrics"
	"00-newapp-template/internal/pkg/server/middleware"
	"00-newapp-template/pkg/acme"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (s *Server) shutdown(w http.ResponseWriter, r *http.Request) {
	s.Log.Debugf("/shutdown called - beginning s shutdown")

	w.Write([]byte("bye felcia"))
	timeout, cancel := context.WithTimeout(s.Context, 5*time.Second)
	err := s.HTTP.Shutdown(timeout)
	if err != nil {
		s.Log.Errorf("s error during shutdown: %+v", err)
	}
	s.Finished()
	cancel()
}

func (s *Server) gophers(w http.ResponseWriter, r *http.Request) {
	endPoint := acme.EndPoint("Gophers")

	s.Metrics.ServerInc(endPoint, metrics.Methods.Service.Get)

	// Check for a cache hit! :- )
	bb, err := s.cacheFetch(r, endPoint)
	if err == nil && len(bb) > 0 {
		w.Write(bb)
		return
	}

	gophers := s.DB.Gophers()
	s.Metrics.DBInc(endPoint, metrics.Methods.DB.Read)

	b, err := json.Marshal(gophers)
	if err != nil {
		s.Log.Errorf("error marshaling gophers: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	s.cacheStore(r, w, endPoint, b)
	w.Write(b)
}
func (s *Server) gopher(w http.ResponseWriter, r *http.Request) {
	endPoint := acme.EndPoint("Gopher")
	s.Metrics.ServerInc(endPoint, metrics.Methods.Service.Get)

	// Check for a cache hit! :- )
	bb, err := s.cacheFetch(r, endPoint)
	if err == nil && len(bb) > 0 {
		w.Write(bb)
		return
	}

	gopherID := middleware.GopherID(r)
	gophers := s.DB.Gophers()
	s.Metrics.DBInc(endPoint, metrics.Methods.DB.Read)
	for _, gopher := range gophers {
		if string(gopher.ID) == gopherID {
			b, err := json.Marshal(gopher)
			if err != nil {
				s.Log.Errorf("error marshaling gopher: %+v", err)
				break
			}

			s.cacheStore(r, w, endPoint, b)
			w.Write(b)
			return
		}
	}

	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}
func (s *Server) updateGopher(w http.ResponseWriter, r *http.Request) {
	endPoint := acme.EndPoint("Gopher")
	s.Metrics.ServerInc(endPoint, metrics.Methods.Service.Get)

	gopher := acme.Gopher{
		ID:          json.Number(middleware.GopherID(r)),
		Name:        middleware.GopherName(r),
		Description: middleware.GopherDescription(r),
	}

	s.DB.UpdateGopher(gopher)
	s.Metrics.DBInc(endPoint, "update")
	s.cacheClear(r, endPoint)

	s.gopher(w, r)

	return
}
func (s *Server) deleteGopher(w http.ResponseWriter, r *http.Request) {
	endPoint := acme.EndPoint("Gopher")
	s.Metrics.ServerInc(endPoint, metrics.Methods.Service.Delete)

	// Deleting a Gopher impacts Gophers cached response, so clear cache for things too!
	s.cacheClear(r, endPoint)
	s.cacheClear(r, acme.EndPoint("Gophers"))

	gopherID := middleware.GopherID(r)

	s.DB.DeleteGopher(gopherID)
	s.Metrics.DBInc(endPoint, metrics.Methods.DB.Delete)

	s.gophers(w, r)
	return
}

func (s *Server) things(w http.ResponseWriter, r *http.Request) {
	endPoint := acme.EndPoint("Things")
	s.Metrics.ServerInc(endPoint, metrics.Methods.Service.Get)

	// Check for a cache hit! :- )
	bb, err := s.cacheFetch(r, endPoint)
	if err == nil && len(bb) > 0 {
		w.Write(bb)
		return
	}

	gopherID := middleware.GopherID(r)
	things := s.DB.GopherThings(gopherID)
	s.Metrics.DBInc(endPoint, metrics.Methods.DB.Read)
	bb, err = json.Marshal(things)
	if err != nil {
		s.Log.Errorf("error marshaling things: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	s.cacheStore(r, w, endPoint, bb)

	w.Write(bb)
}
func (s *Server) thing(w http.ResponseWriter, r *http.Request) {
	endPoint := acme.EndPoint("Thing")
	s.Metrics.ServerInc(endPoint, metrics.Methods.Service.Get)

	// Check for a cache hit! :- )
	bb, err := s.cacheFetch(r, endPoint)
	if err == nil && len(bb) > 0 {
		w.Write(bb)
		return
	}

	thingID := middleware.ThingID(r)
	gopherID := middleware.GopherID(r)

	things := s.DB.GopherThings(gopherID)
	s.Metrics.DBInc(endPoint, metrics.Methods.DB.Read)

	for _, thing := range things {
		if string(thing.ID) == thingID {
			bb, err := json.Marshal(thing)
			if err != nil {
				s.Log.Errorf("error marshaling thing: %+v", err)
				return
			}
			s.cacheStore(r, w, endPoint, bb)
			w.Write(bb)
			return
		}
	}
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}
func (s *Server) updateThing(w http.ResponseWriter, r *http.Request) {
	endPoint := acme.EndPoint("Thing")
	s.Metrics.ServerInc(endPoint, metrics.Methods.Service.Update)
	// TODO: Write and implement s.DB.UpdateThing(thing)
	s.Metrics.DBInc(endPoint, metrics.Methods.DB.Update)
	s.cacheClear(r, endPoint)
}
func (s *Server) deleteThing(w http.ResponseWriter, r *http.Request) {
	endPoint := acme.EndPoint("Thing")
	s.Metrics.ServerInc(endPoint, metrics.Methods.Service.Delete)

	s.cacheClear(r, endPoint)
	s.cacheClear(r, acme.EndPoint("Things"))

	gopherID := middleware.GopherID(r)
	thingID := middleware.ThingID(r)

	s.DB.DeleteThing(gopherID, thingID)
	s.Metrics.DBInc(endPoint, metrics.Methods.DB.Delete)

	s.things(w, r)
}

func (s *Server) cacheClear(r *http.Request, endPoint acme.EndPoint) {
	if s.DiskCache == nil {
		return
	}
	s.Metrics.CacheInc(endPoint, metrics.Methods.Cache.Invalidate)

	filename, _ := acme.ToCacheFilename(endPoint, middleware.ContextMap(r))
	filename = fmt.Sprintf("%s/%s", s.DiskCache.CacheFolder, filename)
	s.DiskCache.Clear(filename)
}
func (s *Server) cacheStore(r *http.Request, w http.ResponseWriter, endPoint acme.EndPoint, bb []byte) {
	if s.DiskCache == nil {
		return
	}
	s.Metrics.CacheInc(endPoint, metrics.Methods.Cache.Store)

	filename, _ := acme.ToCacheFilename(endPoint, middleware.ContextMap(r))
	prettyCache := middleware.NewPrettyPrint(w).Prettify(bb)
	s.DiskCache.Store(filename, prettyCache)
}
func (s *Server) cacheFetch(r *http.Request, endPoint acme.EndPoint) (bb []byte, err error) {
	if s.DiskCache == nil {
		return
	}

	filename, _ := acme.ToCacheFilename(endPoint, middleware.ContextMap(r))
	filename = fmt.Sprintf("%s/%s", s.DiskCache.CacheFolder, filename)
	bb, err = s.DiskCache.Fetch(filename)

	if err == nil && len(bb) > 0 {
		s.Metrics.CacheInc(endPoint, metrics.Methods.Cache.Hit)
	} else {
		s.Metrics.CacheInc(endPoint, metrics.Methods.Cache.Miss)
	}

	return
}
