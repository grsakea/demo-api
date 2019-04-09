package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestPingStatusCode(t *testing.T) {
	router := setupRouter()
	w := performRequest(router, "GET", "/ping/jb")

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPingContent(t *testing.T) {
	router := setupRouter()
	w := performRequest(router, "GET", "/ping/jb")

	assert.Contains(t, w.Body.String(), "jb")
}
