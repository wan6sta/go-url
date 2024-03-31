package config

import (
	"flag"
)

type Config struct {
	BaseURL string
	HTTPServer
}

type HTTPServer struct {
	Address string
}

func NewConfig() *Config {
	var HTTPAddress string
	var BaseURL string

	flag.StringVar(&HTTPAddress, "a", "localhost:8080", "http server address")
	flag.StringVar(&BaseURL, "b", "http://localhost:8080", "final url address")

	flag.Parse()

	return &Config{
		BaseURL: BaseURL,
		HTTPServer: HTTPServer{
			Address: HTTPAddress,
		},
	}
}
