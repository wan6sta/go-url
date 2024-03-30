package main

import (
	"github.com/wan6sta/go-url/internal/config"
	"github.com/wan6sta/go-url/internal/server"
)

func main() {
	cfg := config.NewConfig()
	serv := server.NewAppServer(cfg)

	serv.Run()
}
