package handlers

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wan6sta/go-url/internal/config"
	"github.com/wan6sta/go-url/internal/server"
	"io"
	"net/http"
	"strings"
	"testing"
)

const googleURL = "https://www.google.com"

func TestHandlers(t *testing.T) {
	cfg := config.NewConfig()
	serv := server.NewAppServer(cfg)
	ts := serv.TS

	var ID string

	t.Run("[POST] positive test #1", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", bufio.NewReader(strings.NewReader(googleURL)))

		res, err := ts.Client().Do(request)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, http.StatusCreated)

		resBody, err := io.ReadAll(res.Body)
		resSlice := strings.Split(string(resBody), "/")
		resID := resSlice[len(resSlice)-1]

		ID = resID

		require.NoError(t, err)
		assert.True(t, strings.Contains(string(resBody), cfg.BaseURL))
		assert.Equal(t, res.Header.Get("Content-Type"), "text/plain; charset=utf-8")
	})

	t.Run("[GET] positive test #2", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/"+ID, nil)

		res, err := ts.Client().Do(request)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, http.StatusTemporaryRedirect)
		assert.Equal(t, res.Header.Get("Content-Type"), "text/plain; charset=utf-8")
		assert.Equal(t, res.Header.Get("Location"), googleURL)
	})

	t.Run("[GET] negative test #3", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)

		res, err := ts.Client().Do(request)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, http.StatusBadRequest)
		assert.Equal(t, res.Header.Get("Content-Type"), "text/plain; charset=utf-8")
	})

	t.Run("[PUT] negative test #4", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPut, "/", nil)

		res, err := ts.Client().Do(request)
		require.NoError(t, err)
		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, http.StatusBadRequest)
		assert.Equal(t, res.Header.Get("Content-Type"), "text/plain; charset=utf-8")
	})
}
