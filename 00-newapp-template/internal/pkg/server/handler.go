package server

import (
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
	endPoint := acme.ServiceEndPoint("Gophers")

	// Check for a cache hit! :- )
	bb, err := s.cacheFetch(r, endPoint)
	if err == nil && len(bb) > 0 {
		w.Write(bb)
		return
	}

	gophers := s.DB.Gophers()
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
	endPoint := acme.ServiceEndPoint("Gopher")

	// Check for a cache hit! :- )
	bb, err := s.cacheFetch(r, endPoint)
	if err == nil && len(bb) > 0 {
		w.Write(bb)
		return
	}

	gopherID := middleware.GopherID(r)
	gophers := s.DB.Gophers()
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

func (s *Server) things(w http.ResponseWriter, r *http.Request) {
	endPoint := acme.ServiceEndPoint("Things")

	// Check for a cache hit! :- )
	bb, err := s.cacheFetch(r, endPoint)
	if err == nil && len(bb) > 0 {
		w.Write(bb)
		return
	}

	gopherID := middleware.GopherID(r)
	things := s.DB.GopherThings(gopherID)
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
	endPoint := acme.ServiceEndPoint("Thing")

	// Check for a cache hit! :- )
	bb, err := s.cacheFetch(r, endPoint)
	if err == nil && len(bb) > 0 {
		w.Write(bb)
		return
	}

	thingID := middleware.ThingID(r)
	gopherID := middleware.GopherID(r)

	things := s.DB.GopherThings(gopherID)
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

func (s *Server) updateGopher(w http.ResponseWriter, r *http.Request) {
	endPoint := acme.ServiceEndPoint("Gopher")

	gopher := acme.Gopher{
		ID:          json.Number(middleware.GopherID(r)),
		Name:        middleware.GopherName(r),
		Description: middleware.GopherDescription(r),
	}

	s.DB.UpdateGopher(gopher)
	s.cacheClear(r, endPoint)

	s.gopher(w, r)

	return
}
func (s *Server) deleteGopher(w http.ResponseWriter, r *http.Request) {
	epGopher := acme.ServiceEndPoint("Gopher")
	epGophers := acme.ServiceEndPoint("Gophers")
	s.cacheClear(r, epGopher)
	s.cacheClear(r, epGophers)

	gopherID := middleware.GopherID(r)

	s.DB.DeleteGopher(gopherID)
	s.gophers(w, r)
	return

}

func (s *Server) updateThing(w http.ResponseWriter, r *http.Request) {
	epThing := acme.ServiceEndPoint("Thing")
	s.cacheClear(r, epThing)
}
func (s *Server) deleteThing(w http.ResponseWriter, r *http.Request) {
	epThing := acme.ServiceEndPoint("Thing")
	s.cacheClear(r, epThing)

	gopherID := middleware.GopherID(r)
	thingID := middleware.ThingID(r)
	s.DB.DeleteThing(gopherID, thingID)
	s.DB.GopherThings(gopherID)

	s.things(w, r)
}

func (s *Server) cacheClear(r *http.Request, endpoint acme.ServiceEndPoint) {
	if s.DiskCache == nil {
		return
	}
	filename, _ := acme.ToCacheFilename(endpoint, middleware.ContextMap(r))
	filename = fmt.Sprintf("%s/%s", s.DiskCache.CacheFolder, filename)
	s.DiskCache.Clear(filename)

}
func (s *Server) cacheStore(r *http.Request, w http.ResponseWriter, endpoint acme.ServiceEndPoint, bb []byte) {
	if s.DiskCache == nil {
		return
	}

	filename, _ := acme.ToCacheFilename(endpoint, middleware.ContextMap(r))
	// filename = fmt.Sprintf("%s/%s", s.DiskCache.CacheFolder, filename)
	prettyCache := middleware.NewPrettyPrint(w).Prettify(bb)
	s.DiskCache.Store(filename, prettyCache)

}
func (s *Server) cacheFetch(r *http.Request, endpoint acme.ServiceEndPoint) (bb []byte, err error) {
	if s.DiskCache == nil {
		return
	}

	filename, _ := acme.ToCacheFilename(endpoint, middleware.ContextMap(r))
	filename = fmt.Sprintf("%s/%s", s.DiskCache.CacheFolder, filename)
	bb, err = s.DiskCache.Fetch(filename)
	return
}
