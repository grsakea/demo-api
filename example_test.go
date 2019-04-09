package main

import (
	"net/http"
	"net/http/httptest"
	"os"
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

func TestInitDb(t *testing.T) {
	os.Setenv("DB_STRING", ":memory:")
	router := setupRouter()
	w := performRequest(router, "POST", "/initDB")

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestInitDbWhenInvalidDb(t *testing.T) {
	os.Setenv("DB_STRING", "/invalid_path")
	router := setupRouter()
	w := performRequest(router, "POST", "/initDB")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
