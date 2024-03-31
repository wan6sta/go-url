package server

import (
	"github.com/wan6sta/go-url/internal/config"
	httprouter "github.com/wan6sta/go-url/internal/server/http"
	"log/slog"
	"net/http"
	"net/http/httptest"
)

type AppServer struct {
	TS  *httptest.Server
	Cfg *config.Config
	log *slog.Logger
}

func NewAppServer(cfg *config.Config, log *slog.Logger) *AppServer {
	router := httprouter.NewRouter(cfg, log)
	ts := httptest.NewServer(router.R)

	return &AppServer{
		TS:  ts,
		Cfg: cfg,
		log: log,
	}
}

func (s *AppServer) Run() {
	router := httprouter.NewRouter(s.Cfg, s.log)

	s.log.Info("server started", "address", s.Cfg.HTTPServer)

	err := http.ListenAndServe(s.Cfg.HTTPServer.Address, router.R)
	if err != nil {
		s.log.Error("server stopped", err)
	}
}
