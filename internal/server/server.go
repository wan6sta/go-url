package server

import (
	"fmt"
	"github.com/wan6sta/go-url/internal/config"
	httprouter "github.com/wan6sta/go-url/internal/server/http"
	"net/http"
	"net/http/httptest"
)

type AppServer struct {
	TS  *httptest.Server
	Cfg *config.Config
}

func NewAppServer(cfg *config.Config) *AppServer {
	router := httprouter.NewRouter(cfg)
	ts := httptest.NewServer(router.R)

	return &AppServer{
		TS:  ts,
		Cfg: cfg,
	}
}

func (s *AppServer) Run() {
	router := httprouter.NewRouter(s.Cfg)
	err := http.ListenAndServe(fmt.Sprintf(":%s", s.Cfg.HTTPServer.Port), router.R)
	if err != nil {
		panic(err)
	}
}
