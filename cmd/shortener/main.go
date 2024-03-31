package main

import (
	"github.com/wan6sta/go-url/internal/config"
	"github.com/wan6sta/go-url/internal/logger"
	"github.com/wan6sta/go-url/internal/server"
)

func main() {
	cfg := config.NewConfig()
	log := logger.NewLogger().Sl
	serv := server.NewAppServer(cfg, log)

	serv.Run()
}
