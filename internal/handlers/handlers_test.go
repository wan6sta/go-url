package handlers

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wan6sta/go-url/internal/config"
	"github.com/wan6sta/go-url/internal/repositories"
	"github.com/wan6sta/go-url/internal/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const googleURL = "https://www.google.com"

func TestHandlers(t *testing.T) {
	cfg := config.NewConfig()
	s := storage.NewStorage()
	r := repositories.NewRepository(s, cfg)
	h := NewHandlers(r)
	var ID string

	t.Run("[POST] positive test #1", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/", bufio.NewReader(strings.NewReader(googleURL)))
		w := httptest.NewRecorder()
		h.AppHandler(w, request)
		res := w.Result()

		assert.Equal(t, res.StatusCode, http.StatusCreated)
		defer res.Body.Close()

		resBody, err := io.ReadAll(res.Body)
		resSlice := strings.Split(string(resBody), "/")
		resID := resSlice[len(resSlice)-1]

		ID = resID

		require.NoError(t, err)
		assert.True(t, strings.Contains(string(resBody), cfg.BaseURL))
		assert.Equal(t, res.Header.Get("Content-Type"), "text/plain; charset=utf-8")
	})

	t.Run("[GET] positive test #2", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/"+ID, nil)
		w := httptest.NewRecorder()
		h.AppHandler(w, request)
		res := w.Result()

		assert.Equal(t, res.StatusCode, http.StatusTemporaryRedirect)
		defer res.Body.Close()

		assert.Equal(t, res.Header.Get("Content-Type"), "text/plain; charset=utf-8")
		assert.Equal(t, res.Header.Get("Location"), googleURL)
	})

	t.Run("[GET] negative test #3", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		h.AppHandler(w, request)
		res := w.Result()

		assert.Equal(t, res.StatusCode, http.StatusBadRequest)
		defer res.Body.Close()

		assert.Equal(t, res.Header.Get("Content-Type"), "text/plain; charset=utf-8")
	})

	t.Run("[PUT] negative test #4", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPut, "/", nil)
		w := httptest.NewRecorder()
		h.AppHandler(w, request)
		res := w.Result()

		assert.Equal(t, res.StatusCode, http.StatusBadRequest)
		defer res.Body.Close()

		assert.Equal(t, res.Header.Get("Content-Type"), "text/plain; charset=utf-8")
	})
}
