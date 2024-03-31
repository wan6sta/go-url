package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wan6sta/go-url/internal/config"
	"github.com/wan6sta/go-url/internal/logger"
	"github.com/wan6sta/go-url/internal/server"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const googleURL = "https://www.google.com"

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	req.Header.Set("Content-Type", "text/plain")
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestHandlers(t *testing.T) {
	cfg := config.NewConfig()
	log := logger.NewLogger().Sl
	serv := server.NewAppServer(cfg, log)
	ts := serv.TS

	var ID string

	t.Run("[POST] positive test #1", func(t *testing.T) {
		res, resBody := testRequest(t, ts, http.MethodPost, "/", strings.NewReader(googleURL))

		assert.Equal(t, res.StatusCode, http.StatusCreated)

		resSlice := strings.Split(resBody, "/")
		resID := resSlice[len(resSlice)-1]
		ID = resID

		defer res.Body.Close()

		assert.True(t, strings.Contains(string(resBody), cfg.BaseURL))
		assert.Equal(t, res.Header.Get("Content-Type"), "text/plain; charset=utf-8")
	})

	t.Run("[GET] positive test #2", func(t *testing.T) {
		res, _ := testRequest(t, ts, http.MethodGet, "/"+ID, nil)

		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, http.StatusOK)
	})

	t.Run("[GET] negative test #3", func(t *testing.T) {
		res, _ := testRequest(t, ts, http.MethodGet, "/123", nil)

		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, http.StatusBadRequest)
		assert.Equal(t, res.Header.Get("Content-Type"), "text/plain; charset=utf-8")
	})

	t.Run("[PUT] negative test #4", func(t *testing.T) {
		res, _ := testRequest(t, ts, http.MethodPut, "/123", nil)

		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, http.StatusBadRequest)
		assert.Equal(t, res.Header.Get("Content-Type"), "text/plain; charset=utf-8")
	})
}
