package server

import (
	"00-newapp-template/pkg/server/middleware"
	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
)

// EnableDefaultRouter defines routes with middleware for request tracking, logging, param contexts
func (s *Server) EnableDefaultRouter() {
	s.Router.Use(chimiddleware.RequestID)
	s.Router.Use(middleware.NewStructuredLogger(s.Log))
	s.Router.Use(chimiddleware.Recoverer)

	s.Router.Route("/", func(r chi.Router) {

		r.Use(middleware.InitialCtx)
		r.Use(middleware.PrettyResponseCtx)

		r.Get("/shutdown", s.Shutdown) // Anyone can Shutdown s - try it by visiting http://localhost:10201/shutdown
		r.Get("/gophers", s.Gophers)   // Anyone can get all Gophers etc..
		r.Put("/gophers", s.AddGopher)

		r.Route("/gopher", func(r chi.Router) {
			r.Route("/{GopherID}", func(r chi.Router) {
				r.Use(middleware.GopherCtx)
				r.Get("/", s.Gopher)
				r.Post("/", s.UpdateGopher) // Update/Delete required IsAuthenticated() true.
				r.Delete("/", s.DeleteGopher)

				// Things doesn't a ThingID and therefore doesn't have a ThingCtx
				r.Get("/things", s.Things)

				r.Route("/thing/{ThingID}", func(r chi.Router) {
					r.Use(middleware.ThingCtx) // Requires IsAuthenticated() true.
					r.Get("/", s.Thing)
					r.Post("/", s.UpdateThing)
					r.Delete("/", s.DeleteThing)
				})
			})
		})
	})
}
