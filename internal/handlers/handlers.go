package handlers

import (
	"bytes"
	"encoding/json"
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

type CreateURLRequest struct {
	URL string `json:"url"`
}

type CreateURLResponse struct {
	Result string `json:"result"`
}

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

func (h *Handlers) CreateURLJSONHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	const op = "handlers.CreateURLJSONHandler"
	var cr CreateURLRequest

	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		h.log.Info("cannot parse req.Body")
		http.Error(w, ErrAppInternal.Error(), http.StatusBadRequest)
	}

	err = json.Unmarshal(buf.Bytes(), &cr)
	if err != nil {
		h.log.Info("cannot unmarshal req.Body")
		http.Error(w, ErrAppInternal.Error(), http.StatusBadRequest)
	}

	id, err := h.r.CreateURL(cr.URL)
	if err != nil {
		http.Error(w, ErrAppInternal.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	result := CreateURLResponse{
		Result: id,
	}

	res, err := json.Marshal(result)
	if err != nil {
		h.log.Info("cannot marshal result")
		http.Error(w, ErrAppInternal.Error(), http.StatusBadRequest)
	}

	_, err = w.Write(res)
	if err != nil {
		http.Error(w, ErrAppInternal.Error(), http.StatusBadRequest)
		return
	}
}

func (h *Handlers) NotAllowedHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}
