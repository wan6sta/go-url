package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	BaseURL     string `yaml:"base_url" env:"BASE_URL"`
	StoragePath string `env:"FILE_STORAGE_PATH"`
	HTTPServer
}

type HTTPServer struct {
	Address string `yaml:"server_address" env:"SERVER_ADDRESS"`
}

func NewConfig() *Config {
	var cfg Config

	var HTTPAddress string
	var BaseURL string
	var StoragePath string

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err)
	}

	flag.StringVar(&HTTPAddress, "a", "localhost:8080", "http server address")
	flag.StringVar(&BaseURL, "b", "http://localhost:8080", "final url address")
	flag.StringVar(&StoragePath, "f", "tmp/short-url-db.json", "storage path")

	flag.Parse()

	if cfg.HTTPServer.Address == "" {
		cfg.HTTPServer.Address = HTTPAddress
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = BaseURL
	}

	if cfg.StoragePath == "" {
		cfg.StoragePath = StoragePath
	}

	return &cfg
}
