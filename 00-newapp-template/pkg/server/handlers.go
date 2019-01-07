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

	metricType := metrics.EndPoints.Gophers
	s.Metrics.ServerInc(metricType, metrics.Methods.Service.Get)

	// Check for a cache hit! :- )
	bb, err := s.cacheFetch(r, endPoint, metricType)
	if err == nil && len(bb) > 0 {
		_, _ = w.Write(bb)
		return
	}

	gophers := s.DB.Gophers()
	s.Metrics.DBInc(metricType, metrics.Methods.DB.Read)

	b, err := json.Marshal(gophers)
	if err != nil {
		s.Log.Errorf("error marshaling gophers: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	s.cacheStore(r, w, endPoint, metricType, b)
	_, _ = w.Write(b)
}
func (s *Server) gopher(w http.ResponseWriter, r *http.Request) {
	endPoint := acme.EndPoints.Gopher
	metricType := metrics.EndPoints.Gopher
	s.Metrics.ServerInc(metricType, metrics.Methods.Service.Get)

	// Check for a cache hit! :- )
	bb, err := s.cacheFetch(r, endPoint, metricType)
	if err == nil && len(bb) > 0 {
		_, _ = w.Write(bb)
		return
	}

	gopherID := middleware.GopherID(r)
	gophers := s.DB.Gophers()
	s.Metrics.DBInc(metricType, metrics.Methods.DB.Read)
	for _, gopher := range gophers {
		if string(gopher.ID) == gopherID {
			b, err := json.Marshal(gopher)
			if err != nil {
				s.Log.Errorf("error marshaling gopher: %+v", err)
				break
			}

			s.cacheStore(r, w, endPoint, metricType, b)
			_, _ = w.Write(b)
			return
		}
	}

	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
func (s *Server) updateGopher(w http.ResponseWriter, r *http.Request) {
	// Metrics!
	metricType := metrics.EndPoints.Gopher
	s.Metrics.ServerInc(metricType, metrics.Methods.Service.Update)

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
	s.Metrics.DBInc(metricType, metrics.Methods.DB.Update)

	s.cacheClear(r, acme.EndPoints.Gopher, metricType)
	s.cacheClear(r, acme.EndPoints.Gophers, metricType)

	s.gopher(w, r)

	return
}

func (s *Server) addGopher(w http.ResponseWriter, r *http.Request) {
	// Metrics!
	metricType := metrics.EndPoints.Gopher
	methodType := metrics.Methods.Service.Add

	s.Metrics.ServerInc(metricType, methodType)

	var g acme.Gopher
	err := json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	gopher := acme.Gopher{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
	}

	s.DB.AddGopher(gopher)
	// Metrics!
	s.Metrics.DBInc(metricType, metrics.Methods.DB.Insert)

	s.cacheClear(r, acme.EndPoints.Gopher, metricType)
	s.cacheClear(r, acme.EndPoints.Gophers, metricType)

	b, _ := json.Marshal(gopher)
	_, _ = w.Write(b)

	return
}


func (s *Server) deleteGopher(w http.ResponseWriter, r *http.Request) {
	endPoint := acme.EndPoints.Gopher
	metricType := metrics.EndPoints.Gopher

	s.Metrics.ServerInc(metricType, metrics.Methods.Service.Delete)

	// Deleting a Gopher impacts Gophers cached response, so clear cache for things too!
	s.cacheClear(r, endPoint, metricType)
	s.cacheClear(r, acme.EndPoints.Gophers, metricType)

	gopherID := middleware.GopherID(r)

	s.DB.DeleteGopher(gopherID)
	s.Metrics.DBInc(metricType, metrics.Methods.DB.Delete)

	s.gophers(w, r)
	return
}
func (s *Server) things(w http.ResponseWriter, r *http.Request) {
	endPoint := acme.EndPoints.Things
	metricType := metrics.EndPoints.Things
	s.Metrics.ServerInc(metricType, metrics.Methods.Service.Get)

	// Check for a cache hit! :- )
	bb, err := s.cacheFetch(r, endPoint, metricType)
	if err == nil && len(bb) > 0 {
		_, _ = w.Write(bb)
		return
	}

	gopherID := middleware.GopherID(r)
	things := s.DB.GopherThings(gopherID)
	s.Metrics.DBInc(metricType, metrics.Methods.DB.Read)
	bb, err = json.Marshal(things)
	if err != nil {
		s.Log.Errorf("error marshaling things: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	s.cacheStore(r, w, endPoint, metricType, bb)

	_, _ = w.Write(bb)
}
func (s *Server) thing(w http.ResponseWriter, r *http.Request) {
	endPoint := acme.EndPoints.Thing
	metricType := metrics.EndPoints.Thing
	s.Metrics.ServerInc(metricType, metrics.Methods.Service.Get)

	// Check for a cache hit! :- )
	bb, err := s.cacheFetch(r, endPoint, metricType)
	if err == nil && len(bb) > 0 {
		_, _ = w.Write(bb)
		return
	}

	thingID := middleware.ThingID(r)
	gopherID := middleware.GopherID(r)

	things := s.DB.GopherThings(gopherID)
	s.Metrics.DBInc(metricType, metrics.Methods.DB.Read)

	for _, thing := range things {
		if string(thing.ID) == thingID {
			bb, err := json.Marshal(thing)
			if err != nil {
				s.Log.Errorf("error marshaling thing: %+v", err)
				return
			}
			s.cacheStore(r, w, endPoint, metricType, bb)
			_, _ = w.Write(bb)
			return
		}
	}
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
func (s *Server) updateThing(w http.ResponseWriter, r *http.Request) {
	metricType := metrics.EndPoints.Thing

	s.Metrics.ServerInc(metricType, metrics.Methods.Service.Update)

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
	s.Metrics.DBInc(metricType, metrics.Methods.DB.Update)

	s.cacheClear(r, acme.EndPoints.Thing, metricType)
	s.cacheClear(r, acme.EndPoints.Things, metricType)

	s.thing(w, r)

	return
}
func (s *Server) deleteThing(w http.ResponseWriter, r *http.Request) {
	metricType := metrics.EndPoints.Thing

	s.Metrics.ServerInc(metricType, metrics.Methods.Service.Delete)

	gopherID := middleware.GopherID(r)
	thingID := middleware.ThingID(r)

	s.DB.DeleteThing(gopherID, thingID)
	s.Metrics.DBInc(metricType, metrics.Methods.DB.Delete)

	s.cacheClear(r, acme.EndPoints.Thing, metricType)
	s.cacheClear(r, acme.EndPoints.Things, metricType)

	s.things(w, r)
}