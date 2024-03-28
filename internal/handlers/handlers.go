package handlers

import (
	"errors"
	"fmt"
	"github.com/wan6sta/go-url/internal/storage"
	"io"
	"net/http"
	"strings"
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
	r AppRepos
}

func NewHandlers(r AppRepos) *Handlers {
	return &Handlers{r: r}
}

func (h *Handlers) AppHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, ErrAppBadRequest.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if r.Method == http.MethodGet {
		id := strings.TrimPrefix(r.URL.Path, "/")

		URL, err := h.r.GetURL(id)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				http.Error(w, "URL не найден", http.StatusBadRequest)
				return
			}

			http.Error(w, ErrAppBadRequest.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(URL)

		w.Header().Set("Location", URL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}

	if r.Method == http.MethodPost {
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
}
