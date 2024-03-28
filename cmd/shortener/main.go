package main

import (
	"fmt"
	"github.com/wan6sta/go-url/internal/config"
	"github.com/wan6sta/go-url/internal/handlers"
	"github.com/wan6sta/go-url/internal/repositories"
	"github.com/wan6sta/go-url/internal/server"
	"github.com/wan6sta/go-url/internal/storage"
	"net/http"
)

func main() {
	cfg := config.NewConfig()
	s := storage.NewStorage()
	r := repositories.NewRepository(s, cfg)
	h := handlers.NewHandlers(r)
	mux := server.NewAppServerMux(h)

	err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.HTTPServer.Port), mux.AppMux)
	if err != nil {
		panic(err)
	}
}
