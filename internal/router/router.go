package router

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/novaru/billing-service/internal/app/handler"
)

type Router struct {
	handlers *handler.Handlers
}

func New(handlers *handler.Handlers) *Router {
	return &Router{handlers: handlers}
}

func (rt *Router) Setup() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
	}))

	r.Mount("/api/v1", rt.apiRoutes())

	return r
}
