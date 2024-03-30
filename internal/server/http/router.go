package httprouter

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/wan6sta/go-url/internal/config"
	"github.com/wan6sta/go-url/internal/handlers"
	"github.com/wan6sta/go-url/internal/repositories"
	"github.com/wan6sta/go-url/internal/storage"
)

type Router struct {
	R chi.Router
}

func NewRouter(cfg *config.Config) *Router {
	s := storage.NewStorage()
	repos := repositories.NewRepository(s, cfg)
	h := handlers.NewHandlers(repos)
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/", h.CreateUrlHandler)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.GetUrlHandler)
	})

	return &Router{
		R: r,
	}
}
