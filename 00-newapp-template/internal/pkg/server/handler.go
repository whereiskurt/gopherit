package server

import (
	"00-newapp-template/internal/pkg/server/middleware"
	"00-newapp-template/pkg/acme"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

func (server *Server) shutdown(w http.ResponseWriter, r *http.Request) {
	server.Log.Debugf("/shutdown called - beginning server shutdown")

	w.Write([]byte("bye felcia"))
	timeout, cancel := context.WithTimeout(server.Context, 5*time.Second)
	err := server.HTTP.Shutdown(timeout)
	if err != nil {
		server.Log.Errorf("server error during shutdown: %+v", err)
	}
	server.Finished()
	cancel()
}

func (server *Server) gophers(w http.ResponseWriter, r *http.Request) {

	gophers := server.DB.Gophers()
	b, err := json.Marshal(gophers)
	if err != nil {
		server.Log.Errorf("error marshing gophers: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write(b)
}
func (server *Server) gopher(w http.ResponseWriter, r *http.Request) {
	gopherID := middleware.GopherID(r)

	for _, gopher := range server.DB.Gophers() {
		if string(gopher.ID) == gopherID {
			b, err := json.Marshal(gopher)
			if err != nil {
				server.Log.Errorf("error marshing gopher: %+v", err)
				return
			}
			w.Write(b)
			return
		}
	}
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (server *Server) things(w http.ResponseWriter, r *http.Request) {
	gopherID := middleware.GopherID(r)

	things := server.DB.GopherThings(gopherID)
	b, err := json.Marshal(things)
	if err != nil {
		server.Log.Errorf("error marshing things: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
func (server *Server) thing(w http.ResponseWriter, r *http.Request) {
	thingID := middleware.ThingID(r)
	gopherID := middleware.GopherID(r)

	things := server.DB.GopherThings(gopherID)
	for _, thing := range things {
		if string(thing.ID) == thingID {
			b, err := json.Marshal(thing)
			if err != nil {
				server.Log.Errorf("error marshaling thing: %+v", err)
				return
			}
			w.Write(b)
			return
		}
	}
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (server *Server) updateGopher(w http.ResponseWriter, r *http.Request) {
	gopher := acme.Gopher{
		ID:          json.Number(middleware.GopherID(r)),
		Name:        middleware.GopherName(r),
		Description: middleware.GopherDescription(r),
	}

	server.DB.UpdateGopher(gopher)
	server.gopher(w, r)

	return
}
func (server *Server) deleteGopher(w http.ResponseWriter, r *http.Request) {
	gopherID := middleware.GopherID(r)

	server.DB.DeleteGopher(gopherID)
	server.gophers(w, r)
	return

}

func (server *Server) updateThing(w http.ResponseWriter, r *http.Request) {}
func (server *Server) deleteThing(w http.ResponseWriter, r *http.Request) {
	gopherID := middleware.GopherID(r)
	thingID := middleware.ThingID(r)

	server.DB.DeleteThing(gopherID, thingID)
	server.DB.GopherThings(gopherID)

	server.things(w, r)
	return
}
