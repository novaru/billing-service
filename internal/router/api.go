package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/novaru/billing-service/internal/middleware"
)

func (rt *Router) apiRoutes() chi.Router {
	r := chi.NewRouter()

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", rt.handlers.User.Login)
		r.Post("/register", rt.handlers.User.Create)
	})

	r.Route("/plans", func(r chi.Router) {
		r.Get("/", rt.handlers.Plan.FindAll)
		r.Get("/{slug}", rt.handlers.Plan.FindBySlug)
		r.Post("/", rt.handlers.Plan.Create)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(rt.config))
		r.Use(middleware.APIKeyAuth())

		r.Route("/users", func(r chi.Router) {
			r.Get("/", rt.handlers.User.FindAll)
			r.Get("/{id}", rt.handlers.User.FindByID)
			r.Post("/", rt.handlers.User.Create)
		})
	})

	return r
}
