package handlers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wan6sta/go-url/internal/config"
	"github.com/wan6sta/go-url/internal/handlers"
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

func testJSONRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	req.Header.Set("Content-Type", "application/json")
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func testJSONAndContentEncodingGzipRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func testJSONAndAcceptEncodingGzipRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
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
	})

	t.Run("[POST] JSON positive test #1", func(t *testing.T) {
		req := handlers.CreateURLRequest{URL: googleURL}

		r, err := json.Marshal(req)
		if err != nil {
			t.Error("cannot marshal req", err)
		}

		res, resBody := testJSONRequest(t, ts, http.MethodPost, "/api/shorten", bytes.NewBuffer(r))

		assert.Equal(t, res.StatusCode, http.StatusCreated)

		var cr handlers.CreateURLResponse

		err = json.Unmarshal([]byte(resBody), &cr)
		if err != nil {
			t.Error("cannot unmarshal req", err)
		}

		resSlice := strings.Split(cr.Result, "/")
		resID := resSlice[len(resSlice)-1]
		ID = resID

		t.Log("result id:", cr.Result)
		t.Log("id:", ID)

		defer res.Body.Close()

		assert.True(t, strings.Contains(string(resBody), cfg.BaseURL))
		assert.Equal(t, res.Header.Get("Content-Type"), "application/json")
	})

	t.Run("[GET] JSON positive test #2", func(t *testing.T) {
		res, _ := testJSONRequest(t, ts, http.MethodGet, "/"+ID, nil)

		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, http.StatusOK)
	})

	t.Run("[POST] JSON GZIP positive test #1", func(t *testing.T) {
		req := handlers.CreateURLRequest{URL: strings.Repeat(googleURL, 1)}

		r, err := json.Marshal(req)
		if err != nil {
			t.Error("cannot marshal req", err)
		}

		buf := bytes.NewBuffer(nil)

		zb := gzip.NewWriter(buf)
		_, err = zb.Write([]byte(r))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		res, resBody := testJSONAndContentEncodingGzipRequest(t, ts, http.MethodPost, "/api/shorten", buf)

		assert.Equal(t, res.StatusCode, http.StatusCreated)

		var cr handlers.CreateURLResponse

		err = json.Unmarshal([]byte(resBody), &cr)
		if err != nil {
			t.Error("cannot unmarshal req", err)
		}

		resSlice := strings.Split(cr.Result, "/")
		resID := resSlice[len(resSlice)-1]
		ID = resID

		t.Log("result id:", cr.Result)
		t.Log("id:", ID)

		defer res.Body.Close()

		assert.True(t, strings.Contains(string(resBody), cfg.BaseURL))
		assert.Equal(t, res.Header.Get("Content-Type"), "application/json")
	})

	t.Run("[GET] JSON GZIP positive test #2", func(t *testing.T) {
		res, _ := testJSONAndAcceptEncodingGzipRequest(t, ts, http.MethodGet, "/"+ID, nil)

		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, http.StatusOK)
	})
}
