package server

import (
	"00-newapp-template/internal/pkg/server/middleware"
	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
)

// NewRouter defines routes with middlewares for request tracking, logging, param contexts
func (server *Server) NewRouter() {
	server.Router.Use(chimiddleware.RequestID)
	server.Router.Use(middleware.NewStructuredLogger(server.Log))
	server.Router.Use(chimiddleware.Recoverer)

	server.Router.Route("/", func(r chi.Router) {
		r.Use(middleware.InitialCtx)
		r.Use(middleware.PrettyCtx)

		r.Get("/shutdown", server.shutdown) // Anyone can shutdown server - try it by visiting http://localhost:10201/shutdown
		r.Get("/gophers", server.gophers)   // Anyone can get all gophers

		r.Route("/gopher", func(r chi.Router) {
			r.Route("/{GopherID}", func(r chi.Router) {
				r.Use(middleware.GopherCtx)
				r.Get("/", server.gopher)
				r.Put("/", server.updateGopher) // Update/Delete required IsAuthenticated() true.
				r.Delete("/", server.deleteGopher)

				// Things doesn't a ThingID and therefore doesn't have a ThingCtx
				r.Get("/things", server.things)

				r.Route("/thing/{ThingID}", func(r chi.Router) {
					r.Use(middleware.ThingCtx) // Requires IsAuthenticated() true.
					r.Get("/", server.thing)
					r.Put("/", server.updateThing)
					r.Delete("/", server.deleteThing)
				})
			})
		})
	})
}
