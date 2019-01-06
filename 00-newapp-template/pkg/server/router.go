package server

import (
	"00-newapp-template/pkg/server/middleware"
	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
)

// EnableDefaultRouter defines routes with middleware for request tracking, logging, param contexts
func (s *Server) EnableDefaultRouter() {

	s.Handler.Use(chimiddleware.RequestID)
	s.Handler.Use(middleware.NewStructuredLogger(s.Log))
	s.Handler.Use(chimiddleware.Recoverer)

	s.Handler.Route("/", func(r chi.Router) {

		r.Use(middleware.InitialCtx)
		r.Use(middleware.PrettyResponseCtx)

		r.Get("/shutdown", s.shutdown) // Anyone can shutdown s - try it by visiting http://localhost:10201/shutdown
		r.Get("/gophers", s.gophers)   // Anyone can get all gophers

		r.Route("/gopher", func(r chi.Router) {
			r.Route("/{GopherID}", func(r chi.Router) {
				r.Use(middleware.GopherCtx)
				r.Get("/", s.gopher)
				r.Post("/", s.updateGopher) // Update/Delete required IsAuthenticated() true.
				r.Delete("/", s.deleteGopher)

				// Things doesn't a ThingID and therefore doesn't have a ThingCtx
				r.Get("/things", s.things)

				r.Route("/thing/{ThingID}", func(r chi.Router) {
					r.Use(middleware.ThingCtx) // Requires IsAuthenticated() true.
					r.Get("/", s.thing)
					r.Post("/", s.updateThing)
					r.Delete("/", s.deleteThing)
				})
			})
		})
	})
}
