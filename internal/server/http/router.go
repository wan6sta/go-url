package httprouter

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/wan6sta/go-url/internal/config"
	"github.com/wan6sta/go-url/internal/handlers"
	"github.com/wan6sta/go-url/internal/middlewares"
	"github.com/wan6sta/go-url/internal/repositories"
	"github.com/wan6sta/go-url/internal/storage"
	"log/slog"
	"time"
)

type Router struct {
	R chi.Router
}

func NewRouter(cfg *config.Config, log *slog.Logger) *Router {
	s := storage.NewStorage()
	repos := repositories.NewRepository(s, cfg)
	hm := middlewares.NewHTTPMiddlewares(log)
	h := handlers.NewHandlers(repos, log)
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	// r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.AllowContentType("text/plain"))
	r.Use(hm.Log)

	r.Post("/", h.CreateURLHandler)
	r.Get("/{id}", h.GetURLHandler)

	r.MethodNotAllowed(h.NotAllowedHandler)

	return &Router{
		R: r,
	}
}
