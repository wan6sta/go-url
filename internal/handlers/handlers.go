package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/wan6sta/go-url/internal/storage"
	"io"
	"log/slog"
	"net/http"
)

var (
	ErrAppBadRequest = errors.New("bad request")
	ErrAppInternal   = errors.New("internal server error")
)

type AppRepos interface {
	CreateURL(URL string) (string, error)
	GetURL(ID string) (string, error)
}

type Handlers struct {
	r   AppRepos
	log *slog.Logger
}

func NewHandlers(r AppRepos, log *slog.Logger) *Handlers {
	return &Handlers{r: r, log: log}
}

func (h *Handlers) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	ID := chi.URLParam(r, "id")

	URL, err := h.r.GetURL(ID)
	if err != nil {
		if errors.Is(err, storage.ErrURLNotFound) {
			http.Error(w, "URL не найден", http.StatusBadRequest)
			return
		}

		http.Error(w, ErrAppBadRequest.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", URL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handlers) CreateURLHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, ErrAppInternal.Error(), http.StatusBadRequest)
		return
	}

	URL := string(data)

	id, err := h.r.CreateURL(URL)
	if err != nil {
		http.Error(w, ErrAppInternal.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(id))
	if err != nil {
		http.Error(w, ErrAppInternal.Error(), http.StatusBadRequest)
		return
	}
}

func (h *Handlers) NotAllowedHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
}
