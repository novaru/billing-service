package router

import (
	"github.com/go-chi/chi/v5"
)

func (rt *Router) apiRoutes() chi.Router {
	r := chi.NewRouter()

	r.Route("/users", func(r chi.Router) {
		r.Get("/", rt.handlers.User.FindAll)
		r.Get("/{id}", rt.handlers.User.FindByID)
		r.Post("/", rt.handlers.User.Create)
	})

	return r
}
