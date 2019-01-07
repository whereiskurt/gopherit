package server

import (
	"00-newapp-template/pkg/acme"
	"00-newapp-template/pkg/metrics"
	"00-newapp-template/pkg/server/middleware"
	"00-newapp-template/pkg/ui"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

func (s *Server) shutdown(w http.ResponseWriter, r *http.Request) {
	s.Log.Debugf("/shutdown called - beginning s shutdown")

	_, _ = w.Write([]byte(ui.Gopher()))
	_, _ = w.Write([]byte("\n...bye felicia\n"))

	timeout, cancel := context.WithTimeout(s.Context, 5*time.Second)
	err := s.HTTP.Shutdown(timeout)
	if err != nil {
		s.Log.Errorf("server error during shutdown: %+v", err)
	}
	s.Finished()
	cancel()
}
func (s *Server) gophers(w http.ResponseWriter, r *http.Request) {
	endPoint := acme.EndPoints.Gophers

	serviceType := metrics.EndPoints.Gophers
	s.Metrics.ServerInc(serviceType, metrics.Methods.Service.Get)

	// Check for a cache hit! :- )
	bb, err := s.cacheFetch(r, endPoint, serviceType)
	if err == nil && len(bb) > 0 {
		_, _ = w.Write(bb)
		return
	}

	gophers := s.DB.Gophers()
	s.Metrics.DBInc(serviceType, metrics.Methods.DB.Read)

	b, err := json.Marshal(gophers)
	if err != nil {
		s.Log.Errorf("error marshaling gophers: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	s.cacheStore(r, w, endPoint, serviceType, b)
	_, _ = w.Write(b)
}
func (s *Server) gopher(w http.ResponseWriter, r *http.Request) {
	endPoint := acme.EndPoints.Gopher
	serviceType := metrics.EndPoints.Gopher
	s.Metrics.ServerInc(serviceType, metrics.Methods.Service.Get)

	// Check for a cache hit! :- )
	bb, err := s.cacheFetch(r, endPoint, serviceType)
	if err == nil && len(bb) > 0 {
		_, _ = w.Write(bb)
		return
	}

	gopherID := middleware.GopherID(r)
	gophers := s.DB.Gophers()
	s.Metrics.DBInc(serviceType, metrics.Methods.DB.Read)
	for _, gopher := range gophers {
		if string(gopher.ID) == gopherID {
			b, err := json.Marshal(gopher)
			if err != nil {
				s.Log.Errorf("error marshaling gopher: %+v", err)
				break
			}

			s.cacheStore(r, w, endPoint, serviceType, b)
			_, _ = w.Write(b)
			return
		}
	}

	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
func (s *Server) updateGopher(w http.ResponseWriter, r *http.Request) {
	// Metrics!
	serviceType := metrics.EndPoints.Gopher
	methodType := metrics.Methods.Service.Update
	s.Metrics.ServerInc(serviceType, methodType)

	var g acme.Gopher
	err := json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	gopher := acme.Gopher{
		ID:          json.Number(middleware.GopherID(r)),
		Name:        g.Name,
		Description: g.Description,
	}

	s.DB.UpdateGopher(gopher)
	// Metrics!
	s.Metrics.DBInc(serviceType, "update")

	s.cacheClear(r, acme.EndPoints.Gopher, serviceType)
	s.cacheClear(r, acme.EndPoints.Gophers, serviceType)

	s.gopher(w, r)

	return
}
func (s *Server) deleteGopher(w http.ResponseWriter, r *http.Request) {
	endPoint := acme.EndPoints.Gopher
	serviceType := metrics.EndPoints.Gopher

	s.Metrics.ServerInc(serviceType, metrics.Methods.Service.Delete)

	// Deleting a Gopher impacts Gophers cached response, so clear cache for things too!
	s.cacheClear(r, endPoint, serviceType)
	s.cacheClear(r, acme.EndPoints.Gophers, serviceType)

	gopherID := middleware.GopherID(r)

	s.DB.DeleteGopher(gopherID)
	s.Metrics.DBInc(serviceType, metrics.Methods.DB.Delete)

	s.gophers(w, r)
	return
}
func (s *Server) things(w http.ResponseWriter, r *http.Request) {
	endPoint := acme.EndPoints.Things
	serviceType := metrics.EndPoints.Things
	s.Metrics.ServerInc(serviceType, metrics.Methods.Service.Get)

	// Check for a cache hit! :- )
	bb, err := s.cacheFetch(r, endPoint, serviceType)
	if err == nil && len(bb) > 0 {
		_, _ = w.Write(bb)
		return
	}

	gopherID := middleware.GopherID(r)
	things := s.DB.GopherThings(gopherID)
	s.Metrics.DBInc(serviceType, metrics.Methods.DB.Read)
	bb, err = json.Marshal(things)
	if err != nil {
		s.Log.Errorf("error marshaling things: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	s.cacheStore(r, w, endPoint, serviceType, bb)

	_, _ = w.Write(bb)
}
func (s *Server) thing(w http.ResponseWriter, r *http.Request) {
	endPoint := acme.EndPoints.Thing
	serviceType := metrics.EndPoints.Thing
	s.Metrics.ServerInc(serviceType, metrics.Methods.Service.Get)

	// Check for a cache hit! :- )
	bb, err := s.cacheFetch(r, endPoint, serviceType)
	if err == nil && len(bb) > 0 {
		_, _ = w.Write(bb)
		return
	}

	thingID := middleware.ThingID(r)
	gopherID := middleware.GopherID(r)

	things := s.DB.GopherThings(gopherID)
	s.Metrics.DBInc(serviceType, metrics.Methods.DB.Read)

	for _, thing := range things {
		if string(thing.ID) == thingID {
			bb, err := json.Marshal(thing)
			if err != nil {
				s.Log.Errorf("error marshaling thing: %+v", err)
				return
			}
			s.cacheStore(r, w, endPoint, serviceType, bb)
			_, _ = w.Write(bb)
			return
		}
	}
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
func (s *Server) updateThing(w http.ResponseWriter, r *http.Request) {
	serviceType := metrics.EndPoints.Thing

	s.Metrics.ServerInc(serviceType, metrics.Methods.Service.Update)

	var t acme.Thing
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	thing := acme.Thing{
		GopherID:    t.GopherID,
		ID:          t.ID,
		Description: t.Description,
		Name:        t.Name,
	}

	s.DB.UpdateThing(thing)
	s.Metrics.DBInc(serviceType, metrics.Methods.DB.Update)

	s.cacheClear(r, acme.EndPoints.Thing, serviceType)
	s.cacheClear(r, acme.EndPoints.Things, serviceType)

	s.thing(w, r)

	return
}
func (s *Server) deleteThing(w http.ResponseWriter, r *http.Request) {
	serviceType := metrics.EndPoints.Thing

	s.Metrics.ServerInc(serviceType, metrics.Methods.Service.Delete)

	gopherID := middleware.GopherID(r)
	thingID := middleware.ThingID(r)

	s.DB.DeleteThing(gopherID, thingID)
	s.Metrics.DBInc(serviceType, metrics.Methods.DB.Delete)

	s.cacheClear(r, acme.EndPoints.Thing, serviceType)
	s.cacheClear(r, acme.EndPoints.Things, serviceType)

	s.things(w, r)
}