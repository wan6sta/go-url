package config

type Config struct {
	BaseUrl string
	HTTPServer
}

type HTTPServer struct {
	Port string
	Host string
}

func NewConfig() *Config {
	return &Config{
		BaseUrl: "http://localhost:8080",
		HTTPServer: HTTPServer{
			Port: "8080",
			Host: "http://localhost",
		},
	}
}
