package server

import (
	"00-newapp-template/pkg/acme"
	"00-newapp-template/pkg/metrics"
	"00-newapp-template/pkg/server/middleware"
	"00-newapp-template/pkg/ui"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

func (s *Server) Shutdown(w http.ResponseWriter, r *http.Request) {
	s.Log.Debugf("/Shutdown called - beginning s Shutdown")

	_, _ = w.Write([]byte(ui.Gopher()))
	_, _ = w.Write([]byte("\n...bye felicia\n"))

	timeout, cancel := context.WithTimeout(s.Context, 5*time.Second)
	err := s.HTTP.Shutdown(timeout)
	if err != nil {
		s.Log.Errorf("server error during Shutdown: %+v", err)
	}
	s.Finished()
	cancel()
}
func (s *Server) Gophers(w http.ResponseWriter, r *http.Request) {
	s.Metrics.ServerInc(metrics.EndPoints.Gophers, metrics.Methods.Service.Get)

	// Check for a cache hit! :- )
	bb, err := s.cacheFetch(r, acme.EndPoints.Gophers, metrics.EndPoints.Gophers)
	if err == nil && len(bb) > 0 {
		_, _ = w.Write(bb)
		return
	}

	gophers := s.DB.Gophers()
	s.Metrics.DBInc(metrics.EndPoints.Gophers, metrics.Methods.DB.Read)

	b, err := json.Marshal(gophers)
	if err != nil {
		s.Log.Errorf("error marshaling Gophers: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	s.cacheStore(w, r, b, acme.EndPoints.Gophers, metrics.EndPoints.Gophers)
	_, _ = w.Write(b)
}
func (s *Server) Gopher(w http.ResponseWriter, r *http.Request) {
	s.Metrics.ServerInc(metrics.EndPoints.Gopher, metrics.Methods.Service.Get)

	bb, err := s.cacheFetch(r, acme.EndPoints.Gopher, metrics.EndPoints.Gopher)
	if err == nil && len(bb) > 0 {
		_, _ = w.Write(bb)
		return
	}

	gopherID := middleware.GopherID(r)
	gophers := s.DB.Gophers()
	s.Metrics.DBInc(metrics.EndPoints.Gopher, metrics.Methods.DB.Read)
	for _, gopher := range gophers {
		if string(gopher.ID) == gopherID {
			b, err := json.Marshal(gopher)
			if err != nil {
				s.Log.Errorf("error marshaling Gopher: %+v", err)
				break
			}

			s.cacheStore(w, r, b, acme.EndPoints.Gopher, metrics.EndPoints.Gopher)
			_, _ = w.Write(b)
			return
		}
	}

	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
func (s *Server) AddGopher(w http.ResponseWriter, r *http.Request) {
	s.Metrics.ServerInc(metrics.EndPoints.Gopher, metrics.Methods.Service.Add)

	var g acme.Gopher
	gopherJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	err = json.Unmarshal(gopherJSON, &g)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	gopher := acme.Gopher{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
	}
	s.DB.AddGopher(gopher)
	s.Metrics.DBInc(metrics.EndPoints.Gopher, metrics.Methods.DB.Insert)

	s.cacheClear(r, acme.EndPoints.Gopher, metrics.EndPoints.Gopher)
	s.cacheClear(r, acme.EndPoints.Gophers, metrics.EndPoints.Gophers)

	_, _ = w.Write(gopherJSON)

	return
}
func (s *Server) UpdateGopher(w http.ResponseWriter, r *http.Request) {
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

	s.Gopher(w, r)

	return
}
func (s *Server) DeleteGopher(w http.ResponseWriter, r *http.Request) {
	endPoint := acme.EndPoints.Gopher
	metricType := metrics.EndPoints.Gopher

	s.Metrics.ServerInc(metricType, metrics.Methods.Service.Delete)

	// Deleting a Gopher impacts Gophers cached response, so clear cache for Things too!
	s.cacheClear(r, endPoint, metricType)
	s.cacheClear(r, acme.EndPoints.Gophers, metricType)

	gopherID := middleware.GopherID(r)

	s.DB.DeleteGopher(gopherID)
	s.Metrics.DBInc(metricType, metrics.Methods.DB.Delete)

	s.Gophers(w, r)
	return
}
func (s *Server) Things(w http.ResponseWriter, r *http.Request) {
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
		s.Log.Errorf("error marshaling Things: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	s.cacheStore(w, r, bb, endPoint, metricType)

	_, _ = w.Write(bb)
}
func (s *Server) Thing(w http.ResponseWriter, r *http.Request) {
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
				s.Log.Errorf("error marshaling Thing: %+v", err)
				return
			}
			s.cacheStore(w, r, bb, endPoint, metricType)
			_, _ = w.Write(bb)
			return
		}
	}
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
func (s *Server) AddThing(w http.ResponseWriter, r *http.Request) {
	s.Metrics.ServerInc(metrics.EndPoints.Thing, metrics.Methods.Service.Add)

	var g acme.Thing
	thingJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	err = json.Unmarshal(thingJSON, &g)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	thing := acme.Thing{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
		GopherID:    g.GopherID,
	}
	s.DB.AddThing(thing)
	s.Metrics.DBInc(metrics.EndPoints.Thing, metrics.Methods.DB.Insert)

	s.cacheClear(r, acme.EndPoints.Thing, metrics.EndPoints.Thing)
	s.cacheClear(r, acme.EndPoints.Things, metrics.EndPoints.Things)

	_, _ = w.Write(thingJSON)

	return
}
func (s *Server) UpdateThing(w http.ResponseWriter, r *http.Request) {
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

	s.Thing(w, r)

	return
}
func (s *Server) DeleteThing(w http.ResponseWriter, r *http.Request) {
	metricType := metrics.EndPoints.Thing

	s.Metrics.ServerInc(metricType, metrics.Methods.Service.Delete)

	gopherID := middleware.GopherID(r)
	thingID := middleware.ThingID(r)

	s.DB.DeleteThing(gopherID, thingID)
	s.Metrics.DBInc(metricType, metrics.Methods.DB.Delete)

	s.cacheClear(r, acme.EndPoints.Thing, metricType)
	s.cacheClear(r, acme.EndPoints.Things, metricType)

	s.Things(w, r)
}
