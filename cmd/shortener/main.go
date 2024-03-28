package main

import (
	"fmt"
	"github.com/wan6sta/go-url/internal/app/config"
	"github.com/wan6sta/go-url/internal/app/handlers"
	"github.com/wan6sta/go-url/internal/app/repositories"
	"github.com/wan6sta/go-url/internal/app/server"
	"github.com/wan6sta/go-url/internal/app/storage"
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
