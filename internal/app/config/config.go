package config

type Config struct {
	HTTPServer
}

type HTTPServer struct {
	Port string
	Host string
}

func NewConfig() *Config {
	return &Config{
		HTTPServer: HTTPServer{
			Port: "8080",
			Host: "http://localhost",
		},
	}
}
