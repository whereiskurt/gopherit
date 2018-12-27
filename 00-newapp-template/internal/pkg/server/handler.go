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
	gophers := s.DB.Gophers()
	b, err := json.Marshal(gophers)
	if err != nil {
		s.Log.Errorf("error marshaling gophers: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	s.WriteResponse(r, w, endPoint, b)
}
func (s *Server) gopher(w http.ResponseWriter, r *http.Request) {
	endPoint := acme.ServiceEndPoint("Gopher")
	gopherID := middleware.GopherID(r)

	gophers := s.DB.Gophers()
	for _, gopher := range gophers {
		if string(gopher.ID) == gopherID {
			b, err := json.Marshal(gopher)
			if err != nil {
				s.Log.Errorf("error marshaling gopher: %+v", err)
				break
			}

			s.WriteResponse(r, w, endPoint, b)
			return
		}
	}

	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (s *Server) things(w http.ResponseWriter, r *http.Request) {
	endPoint := acme.ServiceEndPoint("Things")
	gopherID := middleware.GopherID(r)

	things := s.DB.GopherThings(gopherID)
	b, err := json.Marshal(things)
	if err != nil {
		s.Log.Errorf("error marshaling things: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	s.WriteResponse(r, w, endPoint, b)
}
func (s *Server) thing(w http.ResponseWriter, r *http.Request) {
	endPoint := acme.ServiceEndPoint("Thing")
	thingID := middleware.ThingID(r)
	gopherID := middleware.GopherID(r)

	things := s.DB.GopherThings(gopherID)
	for _, thing := range things {
		if string(thing.ID) == thingID {
			b, err := json.Marshal(thing)
			if err != nil {
				s.Log.Errorf("error marshaling thing: %+v", err)
				return
			}
			s.WriteResponse(r, w, endPoint, b)
			return
		}
	}
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (s *Server) updateGopher(w http.ResponseWriter, r *http.Request) {
	gopher := acme.Gopher{
		ID:          json.Number(middleware.GopherID(r)),
		Name:        middleware.GopherName(r),
		Description: middleware.GopherDescription(r),
	}

	s.DB.UpdateGopher(gopher)
	s.gopher(w, r)

	return
}
func (s *Server) deleteGopher(w http.ResponseWriter, r *http.Request) {
	gopherID := middleware.GopherID(r)

	s.DB.DeleteGopher(gopherID)
	s.gophers(w, r)
	return

}

func (s *Server) updateThing(w http.ResponseWriter, r *http.Request) {}
func (s *Server) deleteThing(w http.ResponseWriter, r *http.Request) {
	gopherID := middleware.GopherID(r)
	thingID := middleware.ThingID(r)

	s.DB.DeleteThing(gopherID, thingID)
	s.DB.GopherThings(gopherID)

	s.things(w, r)
	return
}

func (s *Server) WriteResponse(r *http.Request, w http.ResponseWriter, endpoint acme.ServiceEndPoint, b []byte) {
	if s.DiskCache != nil {
		filename, _ := acme.ToCacheFilename(endpoint, middleware.ContextMap(r))
		filename = fmt.Sprintf("%s/%s", s.DiskCache.CacheFolder, filename)
		s.DiskCache.Store(filename, b)
	}
	w.Write(b)
}
