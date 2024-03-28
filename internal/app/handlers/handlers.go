package handlers

import (
	"errors"
	"fmt"
	"github.com/wan6sta/go-url/internal/app/storage"
	"io"
	"net/http"
	"strings"
)

var (
	ErrAppBadRequest = errors.New("bad request")
	ErrAppInternal   = errors.New("internal server error")
)

type AppRepos interface {
	CreateUrl(url string) (string, error)
	GetUrl(ID string) (string, error)
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

	if r.Method == http.MethodGet {
		urls := strings.Split(r.URL.String(), "/")
		id := urls[1]

		url, err := h.r.GetUrl(id)
		if err != nil {
			if errors.Is(err, storage.ErrUrlNotFound) {
				http.Error(w, "URL не найден", http.StatusBadRequest)
				return
			}

			http.Error(w, ErrAppBadRequest.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(url)

		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}

	if r.Method == http.MethodPost {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, ErrAppInternal.Error(), http.StatusBadRequest)
			return
		}

		url := string(data)

		id, err := h.r.CreateUrl(url)
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
