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

func TestAddUserControllerWithoutParams(t *testing.T) {
	os.Remove("test.db")
	os.Setenv("DB_STRING", "test.db")
	router := setupRouter()
	performRequest(router, "POST", "/initDB")
	w := performRequest(router, "GET", "/users")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 4, w.Body.Len())
	os.Remove("test.db")
}
