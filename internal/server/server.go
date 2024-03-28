package server

import "net/http"

type AppServerMux struct {
	AppMux *http.ServeMux
}

type Handlers interface {
	AppHandler(w http.ResponseWriter, r *http.Request)
}

func NewAppServerMux(h Handlers) *AppServerMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", h.AppHandler)

	return &AppServerMux{
		AppMux: mux,
	}
}
